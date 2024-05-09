package chat

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

func ChatMode() {
	var prompt string
	var response string
	var userID string

	reader := bufio.NewReader(os.Stdin)

	userID = uuid.New().String()
	fmt.Printf("You: ")
	prompt, _ = reader.ReadString('\n')
	prompt = strings.TrimSpace(prompt)

	for {

		if prompt == "/exit" {
			color.Cyan("Goodbye!")
			os.Exit(0)

		} else if prompt == "" {
			color.Red("No prompt provided")
			prompt = ""

		} else {
			response = chat(userID, prompt)
			color.Green("Model: %s", response)
			prompt = ""

		}

		fmt.Printf("You: ")
		prompt, _ = reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)

	}

}
