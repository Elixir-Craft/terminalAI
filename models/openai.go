package models

import (
	"context"
	"io"
	"log"

	"github.com/Elixir-Craft/terminalAI/configuration"

	openai "github.com/sashabaranov/go-openai"
)

func init() {
	RegisterBackend("openai", func(modelName string) Model {
		return NewOpenAIModel(modelName, &OpenAIConfig{
			// ApiKey:  os.Getenv("OPENAI_API_KEY"),
			ApiKey:  string(configuration.GetConfig("openai-key")),
			BaseURL: "https://api.openai.com/v1",
		})
	})
	RegisterBackend("gpt4all", func(modelName string) Model {
		return NewOpenAIModel(modelName, &OpenAIConfig{
			ApiKey:      "noapikey",
			BaseURL:     "http://localhost:4891/v1",
			MaxTokens:   4096,
			Temperature: 0.7,
			TopP:        0.4,
			Prefix: []openai.ChatCompletionMessage{{
				Role:    openai.ChatMessageRoleSystem,
				Content: `You are an AI assistant who gives a quality response to whatever humans ask of you.`,
			}},
			NoStreaming: true,
		})
	})
}

type OpenAIConfig struct {
	ApiKey            string
	BaseURL           string
	MaxTokens         int
	Temperature, TopP float32
	Prefix            []openai.ChatCompletionMessage
	NoStreaming       bool
}

func NewOpenAIModel(modelName string, cfg *OpenAIConfig) Model {
	clientConfig := openai.DefaultConfig(cfg.ApiKey)
	clientConfig.BaseURL = cfg.BaseURL

	client := openai.NewClientWithConfig(clientConfig)

	return func() Chat {
		var messages []openai.ChatCompletionMessage
		messages = append(messages, cfg.Prefix...)

		return func(ctx context.Context, prompt string) (StreamingOutput, error) {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			})

			req := openai.ChatCompletionRequest{
				Model:       modelName,
				Messages:    messages,
				MaxTokens:   cfg.MaxTokens,
				Temperature: cfg.Temperature,
				TopP:        cfg.TopP,
				Stream:      !cfg.NoStreaming,
			}

			//https://github.com/nomic-ai/gpt4all/issues/1513
			if cfg.NoStreaming {
				res, err := client.CreateChatCompletion(ctx, req)
				if err != nil {
					return nil, err
				}
				msg := res.Choices[0].Message.Content
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: msg,
				})
				c := make(chan string, 1)
				c <- msg
				close(c)
				return c, nil
			}

			stream, err := client.CreateChatCompletionStream(ctx, req)
			if err != nil {
				return nil, err
			}

			c := make(chan string)
			go func() {
				defer close(c)
				msg := ""
				for {
					res, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatal(err)
					}

					delta := res.Choices[0].Delta.Content
					msg += delta
					c <- delta
				}
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: msg,
				})
			}()
			return c, nil
		}
	}
}
