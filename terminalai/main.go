package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/Elixir-Craft/terminalAI/chat"
	"github.com/Elixir-Craft/terminalAI/configuration"
	"github.com/Elixir-Craft/terminalAI/models"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func getPrompt() (string, string, bool) {

	var promptText string

	var input = flag.String("i", "", "Input File Path")
	var output = flag.String("o", "", "Output File Path")
	var prompt = flag.String("p", "", "Prompt")
	var clipBoard = flag.Bool("c", false, "Prompt From Clipboard")
	var version = flag.Bool("v", false, "Version")

	flag.Parse()

	if os.Args[1] == "config" && len(os.Args) > 1 {
		configuration.Config()
		os.Exit(0)
	}

	if os.Args[1] == "chat" {

		color.Cyan("Terminal AI Chat Mode\n\n")

		color.Yellow("Type '/exit' to exit chat mode\n\n")

		chat.ChatMode(getModel())

		return "", "", true
	}

	if *version {
		fmt.Println("Terminal AI v1.0")
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
	// backend := os.Getenv("TERMINAL_AI_BACKEND")
	// model := os.Getenv("TERMINAL_AI_MODEL")

	backend := string(configuration.GetConfig("service"))
	model := string(configuration.GetConfig("model"))

	return models.NewModel(backend, model)
}

func outputResponse(response models.StreamingOutput, output string) {
	var f *os.File
	var err error
	if output == "" || output == "-" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
	_, err = response.WriteTo(f)
	if err != nil {
		log.Fatal("Error writing response to output")
		log.Fatal(err)
	}
	fmt.Println()
}

func main() {
	// fmt.Println("Terminal AI")

	// handle keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	prompt, output, exit := getPrompt()
	if exit {
		return
	}
	model := getModel()

	response, err := model()(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}

	outputResponse(response, output)
}
