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
	RegisterBackend("genai", NewGenAIModel)
}

func NewGenAIModel(modelName string) Model {
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel(modelName)

	return &GenAiModel{model: model}
}

type GenAiModel struct {
	model *genai.GenerativeModel
}

func (model *GenAiModel) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := model.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	response := resp.Candidates[0].Content.Parts[0]
	responseText := fmt.Sprintf("%s", response)

	return responseText, nil
}
