package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func init() {
	RegisterBackend("genai", func(modelName string) Model {
		apiKey := os.Getenv("GOOGLE_AI_API_KEY")

		client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatal(err)
		}

		model := client.GenerativeModel(modelName)
		model.Temperature = genai.Ptr[float32](0.3)

		return func() Chat {
			c := model.StartChat()

			return func(ctx context.Context, prompt string) (string, error) {
				resp, err := c.SendMessage(ctx, genai.Text(prompt))
				if err != nil {
					return "", err
				}

				str := ""
				for _, part := range resp.Candidates[0].Content.Parts {
					str += fmt.Sprintf("%v", part)
				}
				return str, nil
			}
		}
	})
}
