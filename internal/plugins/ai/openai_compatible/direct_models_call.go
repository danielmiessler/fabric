package openai_compatible

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// DirectlyGetModels is used to fetch models directly from the API
// when the standard OpenAI SDK method fails due to a nonstandard format.
// This is useful for providers like Together that return a direct array of models.
func (c *Client) DirectlyGetModels() ([]string, error) {
	url := c.ApiBaseURL.Value
	if url == "" {
		return nil, fmt.Errorf("API base URL not configured")
	}

	// Ensure URL ends with /models
	if !strings.HasSuffix(url, "/models") {
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		url += "models"
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey.Value))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	var body json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	// Try to parse as an object with data field (OpenAI format)
	var openAIFormat struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &openAIFormat); err == nil && len(openAIFormat.Data) > 0 {
		var modelIDs []string
		for _, model := range openAIFormat.Data {
			modelIDs = append(modelIDs, model.ID)
		}
		return modelIDs, nil
	}

	// Try to parse as a direct array (Together format)
	var directArray []struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &directArray); err == nil {
		var modelIDs []string
		for _, model := range directArray {
			modelIDs = append(modelIDs, model.ID)
		}
		return modelIDs, nil
	}

	return nil, fmt.Errorf("unable to parse models response")
}
