package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"terminalAI/chat"
	"terminalAI/models"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
)

func getPrompt() (string, string) {

	var input = flag.String("i", "", "Input File Path")
	var output = flag.String("o", "", "Output File Path")
	var prompt = flag.String("p", "", "Prompt")
	var clipBoard = flag.Bool("c", false, "Prompt From Clipboard")
	var version = flag.Bool("v", false, "Version")

	var chatO = flag.Bool("chat", false, "Chat Mode")

	flag.Parse()

	var promptText string

	if *chatO {
		fmt.Println("Chat Mode")

		chat.ChatMode()

	}

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
			var (
				inputFile []byte
				err       error
			)
			if *input == "-" {
				inputFile, err = io.ReadAll(os.Stdin)
			} else {
				inputFile, err = os.ReadFile(*input)
			}
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

func getModel() models.Model {
	backend := os.Getenv("TERMINAL_AI_BACKEND")
	model := os.Getenv("TERMINAL_AI_MODEL")
	return models.NewModel(backend, model)
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

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	model := getModel()

	prompt, output := getPrompt()

	response, err := model.Generate(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}

	outputResponse(response, output)
}
