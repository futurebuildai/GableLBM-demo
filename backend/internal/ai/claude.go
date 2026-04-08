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
// It supports both a static key and a dynamic KeyStore.
type Client struct {
	staticKey  string
	keyStore   *KeyStore
	httpClient *http.Client
}

// NewClient creates a new Claude API client with a static key.
func NewClient(apiKey string) *Client {
	return &Client{
		staticKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// NewClientWithKeyStore creates a Claude client that reads the key dynamically.
func NewClientWithKeyStore(ks *KeyStore) *Client {
	return &Client{
		keyStore: ks,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// getKey resolves the API key, preferring keystore over static.
func (c *Client) getKey(ctx context.Context) string {
	if c.keyStore != nil {
		return c.keyStore.Get(ctx)
	}
	return c.staticKey
}

// IsConfigured returns true if a key is available.
func (c *Client) IsConfigured(ctx context.Context) bool {
	return c.getKey(ctx) != ""
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

// FreightInvoiceResult holds the extracted freight invoice data.
type FreightInvoiceResult struct {
	TotalAmount   float64 `json:"total_amount"`
	CarrierName   string  `json:"carrier_name"`
	InvoiceNumber string  `json:"invoice_number"`
}

const freightSystemPrompt = `You are a freight invoice extraction assistant for a lumber and building materials dealer.

Your job is to extract the total freight/shipping charge, carrier name, and invoice number from an uploaded freight invoice — this may be a scan, photo, PDF, or digital document.

Return ONLY valid JSON in this exact format:
{"total_amount": 245.50, "carrier_name": "ABC Freight", "invoice_number": "FR-12345"}

Rules:
- total_amount must be a number in dollars (not cents). This is the total freight charge on the invoice.
- carrier_name is the trucking company or freight carrier name
- invoice_number is the carrier's invoice or reference number
- If you cannot determine a field, use an empty string for text fields or 0 for total_amount
- Output ONLY the JSON object, nothing else — no markdown, no explanation`

// ExtractFreightInvoice sends a freight invoice file to Claude for data extraction.
// Returns the extracted total amount, carrier name, and invoice number.
func (c *Client) ExtractFreightInvoice(ctx context.Context, fileBytes []byte, contentType string) (*FreightInvoiceResult, string, error) {
	apiKey := c.getKey(ctx)
	if apiKey == "" {
		return nil, "", fmt.Errorf("no Anthropic API key configured — please enter the freight total manually")
	}

	var content []contentPart

	switch {
	case strings.HasPrefix(contentType, "image/"):
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
				Text: "Extract the freight charge details from this invoice image. Return JSON with total_amount, carrier_name, and invoice_number.",
			},
		}
	case contentType == "application/pdf":
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
				Text: "Extract the freight charge details from this invoice PDF. Return JSON with total_amount, carrier_name, and invoice_number.",
			},
		}
	default:
		return nil, "", fmt.Errorf("unsupported content type for freight invoice: %s", contentType)
	}

	req := messageRequest{
		Model:     anthropicModel,
		MaxTokens: 1024,
		System:    freightSystemPrompt,
		Messages: []messageContent{
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", anthropicAPIURL, bytes.NewReader(body))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("AI request failed — please enter the freight total manually: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr apiError
		if json.Unmarshal(respBody, &apiErr) == nil {
			return nil, "", fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, apiErr.Error.Message)
		}
		return nil, "", fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var msgResp messageResponse
	if err := json.Unmarshal(respBody, &msgResp); err != nil {
		return nil, "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract text from response
	var rawText strings.Builder
	for _, block := range msgResp.Content {
		if block.Type == "text" {
			rawText.WriteString(block.Text)
		}
	}

	raw := rawText.String()

	// Strip markdown code fences if present (e.g. ```json ... ```)
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```") {
		// Remove opening fence (```json or ```)
		if idx := strings.Index(cleaned, "\n"); idx != -1 {
			cleaned = cleaned[idx+1:]
		}
		// Remove closing fence
		if idx := strings.LastIndex(cleaned, "```"); idx != -1 {
			cleaned = cleaned[:idx]
		}
		cleaned = strings.TrimSpace(cleaned)
	}

	// Parse the JSON response
	var result FreightInvoiceResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, raw, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	return &result, raw, nil
}

// systemPrompt instructs Claude how to extract material lists.
const systemPrompt = `You are a material list extraction assistant for a lumber and building materials dealer in the United States.

Your job is to extract structured line items from uploaded material lists — these may be handwritten notes, printed lists, PDFs, spreadsheets, or photos. Lists may be written in ANY language (especially Spanish) or contain misspellings.

For each item you find, output exactly one line in this format:
QUANTITY UOM - DESCRIPTION

Rules:
- ALWAYS translate ALL descriptions into clean, standard American English using proper US lumber/building materials industry terminology
- Correct ALL misspellings, abbreviations, and shorthand — output the proper product name
- Common Spanish translations: "plancha"=sheet/board, "rollos"=rolls, "clavos"=nails, "libras"=lbs, "pulgadas"=inches, "madera"=lumber, "tornillos"=screws, "cemento"=concrete, "tablas"=boards, "vigas"=beams, "tejas"=shingles, "tubos"=pipes, "pegamento"=adhesive, "lija"=sandpaper, "pintura"=paint
- Recognize misspelled brand names: "Tyguek"/"Tyvek"/"tyvec"=Tyvek, "Quikrete"/"quickrete"=Quikrete, "Hardie"/"hardy"=Hardie
- "SOB board" or "SOB" = OSB (Oriented Strand Board) — this is a common handwriting misread
- QUANTITY must be a number (integer or decimal)
- UOM should be one of: pcs, ea, lf, sf, bf, sheets, bags, rolls, bundles, gal, lbs
- DESCRIPTION should include dimensions, species, grade, and any other identifying details in standard US format
- If you cannot determine the quantity, default to 1
- If you cannot determine the UOM, default to pcs
- Output ONLY the extracted lines, nothing else — no headers, no explanations
- Each line item on its own line

Example — handwritten Spanish input:
"12 - 2x4x8 SPF, 5 - 1/2 plancha SOB board, 2 rollos Tyguek, 50 libras 3 puligades clavos"
Expected output:
12 pcs - 2x4x8 SPF Stud
5 sheets - 1/2 OSB 4x8
2 rolls - Tyvek HomeWrap
50 lbs - 3" Common Nails

More examples:
50 pcs - 2x4x8 SPF Stud
25 pcs - 2x6x12 Doug Fir #2
30 sheets - OSB 7/16 4x8
20 bags - Quikrete 80lb`

// ExtractMaterialList sends a file to Claude for material list extraction.
// Supports images (jpeg, png, gif, webp), PDFs, and pre-processed text from spreadsheets.
func (c *Client) ExtractMaterialList(ctx context.Context, fileBytes []byte, contentType string) (string, error) {
	apiKey := c.getKey(ctx)
	if apiKey == "" {
		return "", fmt.Errorf("no Anthropic API key configured")
	}

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
				Text: "Extract all material list items from this image. Translate any non-English text to standard US building materials terminology. Correct all misspellings and brand name errors. Output each item as: QUANTITY UOM - DESCRIPTION",
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
				Text: "Extract all material list items from this PDF. Translate any non-English text to standard US building materials terminology. Correct all misspellings and brand name errors. Output each item as: QUANTITY UOM - DESCRIPTION",
			},
		}
	case contentType == "text/plain" || contentType == "text/csv":
		// Pre-processed text (from spreadsheet conversion)
		content = []contentPart{
			{
				Type: "text",
				Text: "Extract all material list items from this text data. Translate any non-English text to standard US building materials terminology. Correct all misspellings and brand name errors. Output each item as: QUANTITY UOM - DESCRIPTION\n\n" + string(fileBytes),
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
	httpReq.Header.Set("x-api-key", apiKey)
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

// SendMessage sends a text prompt to Claude and returns the text response.
// Used for general-purpose AI tasks like quote editing commands.
func (c *Client) SendMessage(ctx context.Context, systemPromptText, userMessage string) (string, error) {
	apiKey := c.getKey(ctx)
	if apiKey == "" {
		return "", fmt.Errorf("no Anthropic API key configured")
	}

	req := messageRequest{
		Model:     anthropicModel,
		MaxTokens: 4096,
		System:    systemPromptText,
		Messages: []messageContent{
			{
				Role: "user",
				Content: []contentPart{
					{Type: "text", Text: userMessage},
				},
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
	httpReq.Header.Set("x-api-key", apiKey)
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

	var result strings.Builder
	for _, block := range msgResp.Content {
		if block.Type == "text" {
			result.WriteString(block.Text)
		}
	}

	return result.String(), nil
}
