package platform

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Agent struct {
	Name    string
	Request OpenRouterRequest
}

func NewAgent(name string) *Agent {
	return &Agent{Name: name}
}

func (as *Agent) Send(content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("invalid argument: content cannot be empty")
	}

	as.Request.Messages = append(as.Request.Messages, OpenRouterMessage{Role: LlmRoleUser, Content: content})

	openrouterUrl := "https://openrouter.ai/api/v1/chat/completions"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload, err := json.Marshal(as.Request)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		openrouterUrl,
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENROUTER_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error during API call: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openrouter API returned error status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(bodyBytes, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to decode response envelope: %w", err)
	}

	if openRouterResp.Error != nil {
		return "", fmt.Errorf("openrouter API error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("openrouter returned empty choice array")
	}

	jsonResult := openRouterResp.Choices[0].Message.Content

	return jsonResult, nil
}

func (as *Agent) SetInstructions(filePath embed.FS) error {
	path := fmt.Sprintf("%s.md", as.Name)
	content, err := filePath.ReadFile(path)

	if err != nil {
		return fmt.Errorf("failed to read instructions file: %w", err)
	}

	as.Request.Messages = append(as.Request.Messages, OpenRouterMessage{Role: LlmRoleSystem, Content: string(content)})

	return nil
}

func (as *Agent) SetSchema(filePath embed.FS) error {
	path := fmt.Sprintf("%s.json", as.Name)
	content, err := filePath.ReadFile(path)

	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	var rawSchema map[string]any
	if err := json.Unmarshal(content, &rawSchema); err != nil {
		return fmt.Errorf("invalid json in schema file: %w", err)
	}

	as.Request.ResponseFormat.JSONSchema = JSONSchema{Name: as.Name, Strict: true, Schema: rawSchema}
	return nil
}

func (as *Agent) SetJSONRequest(model string) {
	responseFormat := OpenRouterResponseFormat{Type: LLMResponseJSON, JSONSchema: JSONSchema{}}
	as.Request = OpenRouterRequest{Model: model, ResponseFormat: &responseFormat}

}
