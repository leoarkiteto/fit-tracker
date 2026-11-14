// Package infrastructure implements the aiplanning slice outbound ports: the
// Ollama client and the SQLite reads/writes behind planning.
package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"

	"fittracker-api/internal/aiplanning"
)

// defaultEndpoint and defaultModel match the local development setup.
const (
	defaultEndpoint = "http://localhost:11434"
	defaultModel    = "llama3.2"
)

// OllamaClient is the langchaingo-backed implementation of aiplanning.Planner.
type OllamaClient struct {
	llm      llms.Model
	endpoint string
	model    string
}

var _ aiplanning.Planner = (*OllamaClient)(nil)

// NewOllamaClient builds the client, falling back to the default endpoint and
// model when either is empty.
func NewOllamaClient(endpoint, model string) (*OllamaClient, error) {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	if model == "" {
		model = defaultModel
	}

	llm, err := ollama.New(
		ollama.WithServerURL(endpoint),
		ollama.WithModel(model),
	)
	if err != nil {
		return nil, err
	}

	return &OllamaClient{llm: llm, endpoint: endpoint, model: model}, nil
}

// Model is the configured model name.
func (c *OllamaClient) Model() string { return c.model }

// Endpoint is the configured backend URL.
func (c *OllamaClient) Endpoint() string { return c.endpoint }

// Available verifies the Ollama backend is reachable without consuming tokens.
func (c *OllamaClient) Available(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/", nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	return nil
}

// Plan asks the model for a plan and returns its raw answer.
func (c *OllamaClient) Plan(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeGeneric, userPrompt),
	}

	// LangChainGo handles the complexity of the chat completion.
	resp, err := c.llm.GenerateContent(ctx, content, llms.WithJSONMode())
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	return resp.Choices[0].Content, nil
}
