package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/atotto/clipboard"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func String(p genai.Part) string {
	return fmt.Sprintf("%s", p)
}

func configureAPI() (*genai.Client, context.Context, error) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	apiKey := os.Getenv("GOOGLE_AI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx, err
}

func getPrompt() (string, string) {

	var input = flag.String("i", "", "Input File Path")
	var output = flag.String("o", "", "Output File Path")
	var prompt = flag.String("p", "", "Prompt")
	var clipBoard = flag.Bool("c", false, "Prompt From Clipboard")
	var version = flag.Bool("v", false, "Version")

	flag.Parse()

	var promptText string

	if *version {
		fmt.Println("Terminal AI v0.1")
		// Github URL
		fmt.Println("https://github.com/Elixir-Craft/terminalAI")
		os.Exit(0)
	}

	if *input == "" && *prompt == "" && !*clipBoard {
		if len(os.Args) < 2 {
			log.Fatal("No prompt provided")
		}
		promptText = os.Args[1]

	} else {

		if *clipBoard {
			clipContent, err := clipboard.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
			promptText = clipContent + "\n\n" + *prompt
		} else if *input != "" {
			inputFile, err := os.ReadFile(*input)
			if err != nil {
				log.Fatal(err)
			}

			inputFileContent := string(inputFile)
			promptText = inputFileContent + "\n\n" + *prompt

		} else {
			promptText = *prompt
		}

	}

	// fmt.Println(promptText)
	// fmt.Println(*output)
	// os.Exit(0)

	return promptText, *output

}

func outputResponse(response string, output string) {
	if output != "" {
		err := os.WriteFile(output, []byte(response), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(response)
	}

}

func main() {
	// fmt.Println("Terminal AI")
	client, ctx, _ := configureAPI()
	model := client.GenerativeModel("gemini-pro")

	prompt, output := getPrompt()

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		log.Fatal(err)
	}
	response := resp.Candidates[0].Content.Parts[0]

	outputResponse(String(response), output)

}
