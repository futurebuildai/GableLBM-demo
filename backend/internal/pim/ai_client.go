package pim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// TextAIClient calls the Anthropic Messages API for text generation
type TextAIClient struct {
	apiKey string
	model  string
	client *http.Client
}

// NewTextAIClient creates a new Anthropic API client
func NewTextAIClient(apiKey, model string) *TextAIClient {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &TextAIClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

// anthropicRequest is the Anthropic Messages API request body
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse is a simplified Anthropic Messages API response
type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model string `json:"model"`
}

// Generate sends a prompt to the Anthropic API and returns the text response
func (c *TextAIClient) Generate(systemPrompt, userPrompt string, maxTokens int) (string, string, error) {
	if maxTokens == 0 {
		maxTokens = 2048
	}

	reqBody := anthropicRequest{
		Model:     c.model,
		MaxTokens: maxTokens,
		System:    systemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("api call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("anthropic API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(apiResp.Content) == 0 {
		return "", "", fmt.Errorf("empty response from API")
	}

	return apiResp.Content[0].Text, apiResp.Model, nil
}

// ImageAIClient calls the Stability AI API for image generation
type ImageAIClient struct {
	apiKey string
	client *http.Client
}

// NewImageAIClient creates a new Stability AI client
func NewImageAIClient(apiKey string) *ImageAIClient {
	return &ImageAIClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Generate calls Stability AI to generate an image, returns base64 data
func (c *ImageAIClient) Generate(prompt, style string) (string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("prompt", prompt)
	_ = w.WriteField("output_format", "png")
	if style != "" {
		_ = w.WriteField("style_preset", style)
	}
	w.Close()

	req, err := http.NewRequest("POST", "https://api.stability.ai/v2beta/stable-image/generate/core", &buf)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("stability api call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("stability API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Image string `json:"image"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return "data:image/png;base64," + result.Image, nil
}
