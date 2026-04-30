// Package ai provides AI-powered workout planning services using LangChainGo and Ollama.
package ai

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type LangChainClient struct {
	LLM llms.Model
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

	return &LangChainClient{LLM: llm}, nil
}

func (c *LangChainClient) GenerateWorkoutPlan(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeGeneric, userPrompt),
	}

	resp, err := c.LLM.GenerateContent(ctx, content, llms.WithJSONMode())
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	return resp.Choices[0].Content, nil
}
