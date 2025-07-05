package anthropic

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	"golang.org/x/oauth2"
)

const defaultBaseUrl = "https://api.anthropic.com/"
const oauthClientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
const oauthAuthURL = "https://claude.ai/oauth/authorize"
const oauthTokenURL = "https://console.anthropic.com/v1/oauth/token"
const oauthRedirectURL = "https://console.anthropic.com/oauth/code/callback"

const webSearchToolName = "web_search"
const webSearchToolType = "web_search_20250305"
const sourcesHeader = "## Sources"

func NewClient() (ret *Client) {
	vendorName := "Anthropic"
	ret = &Client{}

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	ret.ApiBaseURL = ret.AddSetupQuestion("API Base URL", false)
	ret.ApiBaseURL.Value = defaultBaseUrl
	ret.UseOAuth = ret.AddSetupQuestionBool("Use OAuth login", false)
	ret.ApiKey = ret.PluginBase.AddSetupQuestion("API key", false)

	ret.maxTokens = 4096
	ret.defaultRequiredUserMessage = "Hi"
	ret.models = []string{
		string(anthropic.ModelClaude3_7SonnetLatest), string(anthropic.ModelClaude3_7Sonnet20250219),
		string(anthropic.ModelClaude3_5HaikuLatest), string(anthropic.ModelClaude3_5Haiku20241022),
		string(anthropic.ModelClaude3_5SonnetLatest), string(anthropic.ModelClaude3_5Sonnet20241022),
		string(anthropic.ModelClaude_3_5_Sonnet_20240620), string(anthropic.ModelClaude3OpusLatest),
		string(anthropic.ModelClaude_3_Opus_20240229), string(anthropic.ModelClaude_3_Haiku_20240307),
		string(anthropic.ModelClaudeOpus4_20250514), string(anthropic.ModelClaudeSonnet4_20250514),
	}

	return
}

type Client struct {
	*plugins.PluginBase
	ApiBaseURL *plugins.SetupQuestion
	ApiKey     *plugins.SetupQuestion
	UseOAuth   *plugins.SetupQuestion

	maxTokens                  int
	defaultRequiredUserMessage string
	models                     []string

	client anthropic.Client
}

// oauthTransport is a custom HTTP transport that adds OAuth Bearer token and beta header
type oauthTransport struct {
	client *Client
	base   http.RoundTripper
}

func (t *oauthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	newReq := req.Clone(req.Context())

	// Get current token (may refresh if needed)
	token, err := t.getValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid OAuth token: %w", err)
	}

	// Add OAuth Bearer token
	newReq.Header.Set("Authorization", "Bearer "+token)

	// Add the anthropic-beta header for OAuth
	newReq.Header.Set("anthropic-beta", "oauth-2025-04-20")

	// Set User-Agent to match AI SDK exactly
	newReq.Header.Set("User-Agent", "ai-sdk/anthropic")

	// Remove x-api-key header if present (OAuth doesn't use it)
	newReq.Header.Del("x-api-key")

	return t.base.RoundTrip(newReq)
}

// getValidToken returns a valid access token, refreshing if necessary
func (t *oauthTransport) getValidToken() (string, error) {
	storage, err := common.NewOAuthStorage()
	if err != nil {
		return "", fmt.Errorf("failed to create OAuth storage: %w", err)
	}

	// Load stored token
	token, err := storage.LoadToken("claude")
	if err != nil {
		return "", fmt.Errorf("failed to load stored token: %w", err)
	}
	// If no token exists, run OAuth flow
	if token == nil {
		fmt.Println("No OAuth token found, initiating authentication...")
		newAccessToken, err := runOAuthFlow()
		if err != nil {
			return "", fmt.Errorf("failed to authenticate: %w", err)
		}
		return newAccessToken, nil
	}

	// Check if token needs refresh (5 minute buffer)
	if token.IsExpired(5) {
		fmt.Println("OAuth token expired, refreshing...")
		newAccessToken, err := refreshToken()
		if err != nil {
			// If refresh fails, try re-authentication
			fmt.Println("Token refresh failed, re-authenticating...")
			newAccessToken, err = runOAuthFlow()
			if err != nil {
				return "", fmt.Errorf("failed to refresh or re-authenticate: %w", err)
			}
		}

		return newAccessToken, nil
	}

	return token.AccessToken, nil
}
func (an *Client) Setup() (err error) {
	if err = an.PluginBase.Ask(an.Name); err != nil {
		return
	}

	if plugins.ParseBoolElseFalse(an.UseOAuth.Value) {
		// Check if we have a valid stored token
		storage, err := common.NewOAuthStorage()
		if err != nil {
			return err
		}

		if !storage.HasValidToken("claude", 5) {
			// No valid token, run OAuth flow
			if _, err = runOAuthFlow(); err != nil {
				return err
			}
		}
	}

	err = an.configure()
	return
}

func (an *Client) configure() (err error) {
	opts := []option.RequestOption{}

	if an.ApiBaseURL.Value != "" {
		baseURL := an.ApiBaseURL.Value

		// API keys use v2 endpoint
		// https://github.com/anthropics/anthropic-sdk-go/blob/main/CHANGELOG.md#020-beta1-2025-03-25
		if strings.Contains(baseURL, "-") && !strings.HasSuffix(baseURL, "/v2") {
			baseURL = strings.TrimSuffix(baseURL, "/")
			baseURL = baseURL + "/v2"
		}
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	if plugins.ParseBoolElseFalse(an.UseOAuth.Value) {
		// For OAuth, use Bearer token with custom headers
		// Create custom HTTP client that adds OAuth Bearer token and beta header
		baseTransport := &http.Transport{}
		httpClient := &http.Client{
			Transport: &oauthTransport{
				client: an,
				base:   baseTransport,
			},
		}
		opts = append(opts, option.WithHTTPClient(httpClient))
	} else {
		opts = append(opts, option.WithAPIKey(an.ApiKey.Value))
	}

	an.client = anthropic.NewClient(opts...)
	return
}

func (an *Client) ListModels() (ret []string, err error) {
	return an.models, nil
}

func (an *Client) SendStream(
	msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions, channel chan string,
) (err error) {
	messages := an.toMessages(msgs)
	if len(messages) == 0 {
		close(channel)
		// No messages to send after normalization, consider this a non-error condition for streaming.
		return
	}

	ctx := context.Background()

	stream := an.client.Messages.NewStreaming(ctx, an.buildMessageParams(messages, opts))

	for stream.Next() {
		event := stream.Current()

		// directly send any non-empty delta text
		if event.Delta.Text != "" {
			channel <- event.Delta.Text
		}
	}

	if stream.Err() != nil {
		fmt.Printf("Messages stream error: %v\n", stream.Err())
	}
	close(channel)
	return
}

func (an *Client) buildMessageParams(msgs []anthropic.MessageParam, opts *common.ChatOptions) (
	params anthropic.MessageNewParams) {

	params = anthropic.MessageNewParams{
		Model:       anthropic.Model(opts.Model),
		MaxTokens:   int64(an.maxTokens),
		TopP:        anthropic.Opt(opts.TopP),
		Temperature: anthropic.Opt(opts.Temperature),
		Messages:    msgs,
	}

	// Add Claude Code spoofing system message for OAuth authentication
	if plugins.ParseBoolElseFalse(an.UseOAuth.Value) {
		params.System = []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: "You are Claude Code, Anthropic's official CLI for Claude.",
			},
		}

	}

	if opts.Search {
		// Build the web-search tool definition:
		webTool := anthropic.WebSearchTool20250305Param{
			Name:         webSearchToolName,
			Type:         webSearchToolType,
			CacheControl: anthropic.NewCacheControlEphemeralParam(),
		}

		if opts.SearchLocation != "" {
			webTool.UserLocation.Type = "approximate"
			webTool.UserLocation.Timezone = anthropic.Opt(opts.SearchLocation)
		}

		// Wrap it in the union:
		params.Tools = []anthropic.ToolUnionParam{
			{OfWebSearchTool20250305: &webTool},
		}
	}
	return
}

func (an *Client) Send(ctx context.Context, msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions) (
	ret string, err error) {

	messages := an.toMessages(msgs)
	if len(messages) == 0 {
		// No messages to send after normalization, return empty string and no error.
		return
	}

	var message *anthropic.Message
	if message, err = an.client.Messages.New(ctx, an.buildMessageParams(messages, opts)); err != nil {
		return
	}

	var textParts []string
	var citations []string
	citationMap := make(map[string]bool) // To avoid duplicate citations

	for _, block := range message.Content {
		if block.Type == "text" && block.Text != "" {
			textParts = append(textParts, block.Text)

			// Extract citations from this text block
			for _, citation := range block.Citations {
				if citation.Type == "web_search_result_location" {
					citationKey := citation.URL + "|" + citation.Title
					if !citationMap[citationKey] {
						citationMap[citationKey] = true
						citationText := fmt.Sprintf("- [%s](%s)", citation.Title, citation.URL)
						if citation.CitedText != "" {
							citationText += fmt.Sprintf(" - \"%s\"", citation.CitedText)
						}
						citations = append(citations, citationText)
					}
				}
			}
		}
	}

	var resultBuilder strings.Builder
	resultBuilder.WriteString(strings.Join(textParts, ""))

	// Append citations if any were found
	if len(citations) > 0 {
		resultBuilder.WriteString("\n\n")
		resultBuilder.WriteString(sourcesHeader)
		resultBuilder.WriteString("\n\n")
		resultBuilder.WriteString(strings.Join(citations, "\n"))
	}
	ret = resultBuilder.String()

	return
}

func (an *Client) toMessages(msgs []*chat.ChatCompletionMessage) (ret []anthropic.MessageParam) {
	// Custom normalization for Anthropic:
	// - System messages become the first part of the first user message.
	// - Messages must alternate user/assistant.
	// - Skip empty messages.

	var anthropicMessages []anthropic.MessageParam
	var systemContent string

	// Note: Claude Code spoofing is now handled in buildMessageParams

	isFirstUserMessage := true
	lastRoleWasUser := false

	for _, msg := range msgs {
		if msg.Content == "" {
			continue // Skip empty messages
		}

		switch msg.Role {
		case chat.ChatMessageRoleSystem:
			// Accumulate system content. It will be prepended to the first user message.
			if systemContent != "" {
				systemContent += "\\n" + msg.Content
			} else {
				systemContent = msg.Content
			}
		case chat.ChatMessageRoleUser:
			userContent := msg.Content
			if isFirstUserMessage && systemContent != "" {
				userContent = systemContent + "\\n\\n" + userContent
				isFirstUserMessage = false // System content now consumed
			}
			if lastRoleWasUser {
				// Enforce alternation: add a minimal assistant message if two user messages are consecutive.
				// This shouldn't happen with current chatter.go logic but is a safeguard.
				anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock("Okay.")))
			}
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(userContent)))
			lastRoleWasUser = true
		case chat.ChatMessageRoleAssistant:
			// If the first message is an assistant message, and we have system content,
			// prepend a user message with the system content.
			if isFirstUserMessage && systemContent != "" {
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(systemContent)))
				lastRoleWasUser = true
				isFirstUserMessage = false // System content now consumed
			} else if !lastRoleWasUser && len(anthropicMessages) > 0 {
				// Enforce alternation: add a minimal user message if two assistant messages are consecutive
				// or if an assistant message is first without prior system prompt handling.
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(an.defaultRequiredUserMessage)))
				lastRoleWasUser = true
			}
			anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
			lastRoleWasUser = false
		default:
			// Other roles (like 'meta') are ignored for Anthropic's message structure.
			continue
		}
	}

	// If only system content was provided, create a user message with it.
	if len(anthropicMessages) == 0 && systemContent != "" {
		anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(systemContent)))
	}

	return anthropicMessages
}

func (an *Client) NeedsRawMode(modelName string) bool {
	return false
}

func generatePKCE() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return
}

func openBrowser(url string) {
	commands := [][]string{{"xdg-open", url}, {"open", url}, {"cmd", "/c", "start", url}}
	for _, cmd := range commands {
		if exec.Command(cmd[0], cmd[1:]...).Start() == nil {
			return
		}
	}
}

func runOAuthFlow() (token string, err error) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return
	}

	cfg := oauth2.Config{
		ClientID:    oauthClientID,
		Endpoint:    oauth2.Endpoint{AuthURL: oauthAuthURL, TokenURL: oauthTokenURL},
		RedirectURL: oauthRedirectURL,
		Scopes:      []string{"org:create_api_key", "user:profile", "user:inference"},
	}

	authURL := cfg.AuthCodeURL(verifier,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code", "true"),
		oauth2.SetAuthURLParam("state", verifier),
	)

	fmt.Println("Open the following URL in your browser. Fabric would like to authorize:")
	fmt.Println(authURL)
	openBrowser(authURL)
	fmt.Print("Paste the authorization code here: ")
	var code string
	fmt.Scanln(&code)
	parts := strings.SplitN(code, "#", 2)
	state := verifier
	if len(parts) == 2 {
		state = parts[1]
	}

	// Manual token exchange to match opencode implementation
	tokenReq := map[string]string{
		"code":          parts[0],
		"state":         state,
		"grant_type":    "authorization_code",
		"client_id":     oauthClientID,
		"redirect_uri":  oauthRedirectURL,
		"code_verifier": verifier,
	}

	token, err = exchangeToken(tokenReq)
	return
}

func exchangeToken(params map[string]string) (token string, err error) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		return
	}

	resp, err := http.Post(oauthTokenURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("token exchange failed: %s - %s", resp.Status, string(body))
		return
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	// Save the complete token information
	storage, err := common.NewOAuthStorage()
	if err != nil {
		return result.AccessToken, fmt.Errorf("failed to create OAuth storage: %w", err)
	}

	oauthToken := &common.OAuthToken{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(result.ExpiresIn),
		TokenType:    result.TokenType,
		Scope:        result.Scope,
	}

	if err = storage.SaveToken("claude", oauthToken); err != nil {
		return result.AccessToken, fmt.Errorf("failed to save OAuth token: %w", err)
	}

	token = result.AccessToken
	return
}

// refreshToken refreshes an expired OAuth token using the refresh token
func refreshToken() (string, error) {
	storage, err := common.NewOAuthStorage()
	if err != nil {
		return "", fmt.Errorf("failed to create OAuth storage: %w", err)
	}

	// Load existing token
	token, err := storage.LoadToken("claude")
	if err != nil {
		return "", fmt.Errorf("failed to load stored token: %w", err)
	}
	if token == nil || token.RefreshToken == "" {
		return "", fmt.Errorf("no refresh token available")
	}

	// Prepare refresh request
	refreshReq := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": token.RefreshToken,
		"client_id":     oauthClientID,
	}

	reqBody, err := json.Marshal(refreshReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	// Make refresh request
	resp, err := http.Post(oauthTokenURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token refresh failed: %s - %s", resp.Status, string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// Update stored token
	newToken := &common.OAuthToken{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(result.ExpiresIn),
		TokenType:    result.TokenType,
		Scope:        result.Scope,
	}

	// Use existing refresh token if new one not provided
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = token.RefreshToken
	}

	if err = storage.SaveToken("claude", newToken); err != nil {
		return "", fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return result.AccessToken, nil
}

// runOAuthFlowWithTimeout runs the OAuth flow with an extended timeout for user interaction
func runOAuthFlowWithTimeout() (string, error) {
	// Use a longer timeout context for OAuth flows (10 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Store the context for use in HTTP requests
	return runOAuthFlowWithContext(ctx)
}

// refreshTokenWithTimeout refreshes a token with extended timeout
func refreshTokenWithTimeout() (string, error) {
	// Use a longer timeout context for token refresh (2 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	return refreshTokenWithContext(ctx)
}

// runOAuthFlowWithContext runs OAuth flow with the provided context
func runOAuthFlowWithContext(ctx context.Context) (string, error) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return "", err
	}

	cfg := oauth2.Config{
		ClientID:    oauthClientID,
		Endpoint:    oauth2.Endpoint{AuthURL: oauthAuthURL, TokenURL: oauthTokenURL},
		RedirectURL: oauthRedirectURL,
		Scopes:      []string{"org:create_api_key", "user:profile", "user:inference"},
	}

	authURL := cfg.AuthCodeURL(verifier,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code", "true"),
		oauth2.SetAuthURLParam("state", verifier),
	)

	fmt.Println("Open the following URL in your browser. Fabric would like to authorize:")
	fmt.Println(authURL)
	openBrowser(authURL)
	fmt.Print("Paste the authorization code here: ")

	// Create a channel to receive the user input
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		var code string
		_, err := fmt.Scanln(&code)
		if err != nil {
			errChan <- err
			return
		}
		codeChan <- code
	}()

	// Wait for either user input or context timeout
	var code string
	select {
	case code = <-codeChan:
		// User provided input
	case err := <-errChan:
		return "", fmt.Errorf("failed to read authorization code: %w", err)
	case <-ctx.Done():
		return "", fmt.Errorf("authentication timed out - please try again")
	}

	parts := strings.SplitN(code, "#", 2)
	state := verifier
	if len(parts) == 2 {
		state = parts[1]
	}

	// Manual token exchange to match opencode implementation
	tokenReq := map[string]string{
		"code":          parts[0],
		"state":         state,
		"grant_type":    "authorization_code",
		"client_id":     oauthClientID,
		"redirect_uri":  oauthRedirectURL,
		"code_verifier": verifier,
	}

	return exchangeTokenWithContext(ctx, tokenReq)
}

// refreshTokenWithContext refreshes a token with the provided context
func refreshTokenWithContext(ctx context.Context) (string, error) {
	storage, err := common.NewOAuthStorage()
	if err != nil {
		return "", fmt.Errorf("failed to create OAuth storage: %w", err)
	}

	// Load existing token
	token, err := storage.LoadToken("claude")
	if err != nil {
		return "", fmt.Errorf("failed to load stored token: %w", err)
	}
	if token == nil || token.RefreshToken == "" {
		return "", fmt.Errorf("no refresh token available")
	}

	// Prepare refresh request
	refreshReq := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": token.RefreshToken,
		"client_id":     oauthClientID,
	}

	return exchangeTokenWithContext(ctx, refreshReq)
}

// exchangeTokenWithContext exchanges tokens with the provided context
func exchangeTokenWithContext(ctx context.Context, params map[string]string) (string, error) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", oauthTokenURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange failed: %s - %s", resp.Status, string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	// Save the complete token information
	storage, err := common.NewOAuthStorage()
	if err != nil {
		return "", fmt.Errorf("failed to create OAuth storage: %w", err)
	}

	oauthToken := &common.OAuthToken{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(result.ExpiresIn),
		TokenType:    result.TokenType,
		Scope:        result.Scope,
	}

	// Use existing refresh token if new one not provided (for refresh operations)
	if oauthToken.RefreshToken == "" && params["grant_type"] == "refresh_token" {
		if existingToken, err := storage.LoadToken("claude"); err == nil && existingToken != nil {
			oauthToken.RefreshToken = existingToken.RefreshToken
		}
	}

	if err = storage.SaveToken("claude", oauthToken); err != nil {
		return "", fmt.Errorf("failed to save OAuth token: %w", err)
	}

	return result.AccessToken, nil
}
