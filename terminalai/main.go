package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func configureAPI() (*genai.Client, context.Context, error) {
	err := godotenv.Load()
	if err != nil {
		// Handle error (e.g., file not found)
		panic("Error loading .env file")
	}

	apiKey := os.Getenv("GOOGLE_AI_API_KEY")

	fmt.Println(apiKey)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx, err
}

// defer client.Close()

func main() {
	fmt.Println("Terminal AI")

	client, ctx, _ := configureAPI()

	// fmt.Println(os.Getenv("GOOGLE_AI_API_KEY"))

	model := client.GenerativeModel("gemini-pro")

	// em := client.EmbeddingModel("embedding-001")

	resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.Candidates[0].Content)
	}

	// fmt.Println(resp.Embedding.Values)
}
