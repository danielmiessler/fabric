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
	client.ApiKey = client.AddSetupQuestion("API Key", true)
	return client
}

// return the main content of a webpage in clean, LLM-friendly text.
func (jc *JinaClient) ScrapeURL(url string) (string, error) {
	requestURL := "https://r.jina.ai/" + url
    req, err := http.NewRequest("GET", requestURL, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }

	apiKey := jc.ApiKey.Value

    // Set the Authorization header with the Bearer token
    req.Header.Set("Authorization", "Bearer " + apiKey)

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

// search engine call that returns top-5 results with their URLs and contents, each in clean, LLM-friendly text.
func (jc *JinaClient) ScrapeQuestion(question string) (string, error) {
	url := "https://s.jina.ai/" + question

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}