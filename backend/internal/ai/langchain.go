package ai

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type LangChainClient struct {
	LLM      llms.Model
	Endpoint string
	Model    string
}

func NewLangChainClient(endpoint, model string) (*LangChainClient, error) {
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}

	llm, err := ollama.New(
		ollama.WithServerURL(endpoint),
		ollama.WithModel(model),
	)
	if err != nil {
		return nil, err
	}

	return &LangChainClient{LLM: llm, Endpoint: endpoint, Model: model}, nil
}

// CheckHealth verifies if the Ollama backend is reachable without consuming tokens.
func (c *LangChainClient) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Endpoint+"/", nil)
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

func (c *LangChainClient) GenerateWorkoutPlan(
	ctx context.Context,
	systemPrompt, userPrompt string,
) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeGeneric, userPrompt),
	}

	// LangChainGo handles the complexity of the chat completion
	resp, err := c.LLM.GenerateContent(ctx, content, llms.WithJSONMode())
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	return resp.Choices[0].Content, nil
}
