package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	anthropicAPIURL = "https://api.anthropic.com/v1/messages"
	anthropicModel  = "claude-sonnet-4-5-20250929"
	apiVersion      = "2023-06-01"
)

// Client wraps the Anthropic Messages API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Claude API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// --- Request / Response types ---

type messageRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	System    string           `json:"system"`
	Messages  []messageContent `json:"messages"`
}

type messageContent struct {
	Role    string        `json:"role"`
	Content []contentPart `json:"content"`
}

type contentPart struct {
	Type      string      `json:"type"`
	Text      string      `json:"text,omitempty"`
	Source    *mediaSource `json:"source,omitempty"`
}

type mediaSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type messageResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

type apiError struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// systemPrompt instructs Claude how to extract material lists.
const systemPrompt = `You are a material list extraction assistant for a lumber and building materials dealer.

Your job is to extract structured line items from uploaded material lists — these may be handwritten notes, printed lists, PDFs, spreadsheets, or photos.

For each item you find, output exactly one line in this format:
QUANTITY UOM - DESCRIPTION

Rules:
- QUANTITY must be a number (integer or decimal)
- UOM should be one of: pcs, ea, lf, sf, bf, sheets, bags, rolls, bundles, gal
- DESCRIPTION should include dimensions, species, grade, and any other identifying details
- If you cannot determine the quantity, default to 1
- If you cannot determine the UOM, default to pcs
- Output ONLY the extracted lines, nothing else — no headers, no explanations
- Each line item on its own line
- Preserve the original descriptions as closely as possible while being clear

Example output:
50 pcs - 2x4x8 SPF Stud
25 pcs - 2x6x12 Doug Fir #2
30 sheets - OSB 7/16 4x8
20 bags - Quikrete 80lb`

// ExtractMaterialList sends a file to Claude for material list extraction.
// Supports images (jpeg, png, gif, webp), PDFs, and pre-processed text from spreadsheets.
func (c *Client) ExtractMaterialList(ctx context.Context, fileBytes []byte, contentType string) (string, error) {
	var content []contentPart

	switch {
	case strings.HasPrefix(contentType, "image/"):
		// Image content block
		content = []contentPart{
			{
				Type: "image",
				Source: &mediaSource{
					Type:      "base64",
					MediaType: contentType,
					Data:      base64.StdEncoding.EncodeToString(fileBytes),
				},
			},
			{
				Type: "text",
				Text: "Extract all material list items from this image. Output each item as: QUANTITY UOM - DESCRIPTION",
			},
		}
	case contentType == "application/pdf":
		// PDF document content block
		content = []contentPart{
			{
				Type: "document",
				Source: &mediaSource{
					Type:      "base64",
					MediaType: "application/pdf",
					Data:      base64.StdEncoding.EncodeToString(fileBytes),
				},
			},
			{
				Type: "text",
				Text: "Extract all material list items from this PDF. Output each item as: QUANTITY UOM - DESCRIPTION",
			},
		}
	case contentType == "text/plain" || contentType == "text/csv":
		// Pre-processed text (from spreadsheet conversion)
		content = []contentPart{
			{
				Type: "text",
				Text: "Extract all material list items from this text data. Output each item as: QUANTITY UOM - DESCRIPTION\n\n" + string(fileBytes),
			},
		}
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	req := messageRequest{
		Model:     anthropicModel,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages: []messageContent{
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", anthropicAPIURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr apiError
		if json.Unmarshal(respBody, &apiErr) == nil {
			return "", fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, apiErr.Error.Message)
		}
		return "", fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var msgResp messageResponse
	if err := json.Unmarshal(respBody, &msgResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract text from response
	var result strings.Builder
	for _, block := range msgResp.Content {
		if block.Type == "text" {
			result.WriteString(block.Text)
		}
	}

	return result.String(), nil
}
