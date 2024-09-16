package jina

// see https://jina.ai for more information

import (
	"fmt"
	"io"
	"net/http"

	"github.com/danielmiessler/fabric/common"
)

type JinaClient struct {
	*common.Configurable
	ApiKey *common.SetupQuestion
}

func NewJinaClient() *JinaClient {

	label := "Jina AI"

	client := &JinaClient{
		Configurable: &common.Configurable{
			Label: label,
			EnvNamePrefix: common.BuildEnvVariablePrefix(label),
		},
	}
	client.ApiKey = client.AddSetupQuestion("API Key", false)
    return client
}

// return the main content of a webpage in clean, LLM-friendly text.
func (jc *JinaClient) ScrapeURL(url string) (string, error) {
	requestURL := "https://r.jina.ai/" + url
    req, err := http.NewRequest("GET", requestURL, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

	// if api keys exist, set the header
	if apiKey := jc.ApiKey.Value; apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+apiKey)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    return string(body), nil
}

func (jc *JinaClient) ScrapeQuestion(question string) (string, error) {
    requestURL := "https://s.jina.ai/" + question
    req, err := http.NewRequest("GET", requestURL, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

    // if api keys exist, set the header
    if apiKey := jc.ApiKey.Value; apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+apiKey)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    return string(body), nil
}