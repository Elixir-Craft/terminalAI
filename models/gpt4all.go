package models

import (
	"context"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

func init() {
	RegisterBackend("gpt4all", NewGPT4AllModel)
}

func NewGPT4AllModel(modelName string) Model {
	config := openai.DefaultConfig("noapikey")
	config.BaseURL = "http://localhost:4891/v1"

	client := openai.NewClientWithConfig(config)

	return &GPT4AllModel{
		client:    client,
		modelName: modelName,
	}
}

type GPT4AllModel struct {
	client    *openai.Client
	modelName string
}

func (model *GPT4AllModel) Generate(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model:       model.modelName,
		MaxTokens:   4096,
		Temperature: 0.7,
		TopP:        0.4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: `You are an AI assistant who gives a quality response to whatever humans ask of you.`,
			},
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
