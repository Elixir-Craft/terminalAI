package models

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func init() {
	RegisterBackend("openai", NewOpenAIModel)
}

func NewOpenAIModel(modelName string) Model {
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)

	return &OpenAiModel{
		client:    client,
		modelName: modelName,
	}
}

type OpenAiModel struct {
	client    *openai.Client
	modelName string
}

func (model *OpenAiModel) Generate(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: model.modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	resp, err := model.client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Choices[0].Message.Content, nil
}
