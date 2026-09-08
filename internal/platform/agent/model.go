package platform

type OpenRouterRequest struct {
	Model          string                    `json:"model"`
	Messages       []OpenRouterMessage       `json:"messages"`
	ResponseFormat *OpenRouterResponseFormat `json:"response_format,omitempty"`
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponseFormat struct {
	Type       string     `json:"type"`
	JSONSchema JSONSchema `json:"json_schema"`
}

type JSONSchema struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema map[string]any `json:"schema"`
}

type OpenRouterResponse struct {
	ID      string             `json:"id"`
	Choices []OpenRouterChoice `json:"choices"`
	Error   *OpenRouterError   `json:"error,omitempty"`
}

type OpenRouterChoice struct {
	Message      OpenRouterMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

type OpenRouterError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	LlmRoleSystem string = "system"
	LlmRoleUser   string = "user"
)

var (
	LLMFreeModel                 string = "openrouter/free"
	LLMPaidGeminiFlashThreeSeven string = "google/gemini-3.7-flash"
)

var (
	LLMResponseJSON string = "json_schema"
)
