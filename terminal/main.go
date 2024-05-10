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
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func getPrompt() (string, string, bool) {

	var input = flag.String("i", "", "Input File Path")
	var output = flag.String("o", "", "Output File Path")
	var prompt = flag.String("p", "", "Prompt")
	var clipBoard = flag.Bool("c", false, "Prompt From Clipboard")
	var version = flag.Bool("v", false, "Version")

	var chatMode = flag.Bool("chat", false, "Chat Mode")

	flag.Parse()

	var promptText string

	if *chatMode {

		color.Cyan("Terminal AI Chat Mode\n\n")

		// Instructions
		color.Yellow("Type '/exit' to exit chat mode\n\n")

		chat.ChatMode(getModel())

		return "", "", true
	}

	if *version {
		fmt.Println("Terminal AI v0.1")
		// Github URL
		fmt.Println("https://github.com/Elixir-Craft/terminalAI")
		return "", "", true
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

	return promptText, *output, false

}

func getModel() models.Model {
	backend := os.Getenv("TERMINAL_AI_BACKEND")
	model := os.Getenv("TERMINAL_AI_MODEL")
	return models.NewModel(backend, model)
}

func outputResponse(response io.Reader, output string) {
	if output == "" || output == "-" {
		_, err := io.Copy(os.Stdout, response)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, response)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// fmt.Println("Terminal AI")

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	model := getModel()

	prompt, output, exit := getPrompt()
	if exit {
		return
	}

	response, err := model()(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}

	outputResponse(response, output)
}
