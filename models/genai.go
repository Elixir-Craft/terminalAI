package models

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
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

			return func(ctx context.Context, prompt string) (io.Reader, error) {
				stream := c.SendMessageStream(ctx, genai.Text(prompt))
				if err != nil {
					return nil, err
				}

				r, w := io.Pipe()
				go func() {
					defer w.Close()
					for {
						res, err := stream.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							log.Fatal(err)
						}

						str := ""
						for _, part := range res.Candidates[0].Content.Parts {
							str += fmt.Sprintf("%v", part)
						}
						_, err = w.Write([]byte(str))
						if err != nil {
							log.Fatal(err)
						}
					}
				}()
				return r, nil
			}
		}
	})
}
