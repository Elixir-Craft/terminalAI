package chat

import (
	"bufio"
	"fmt"
	"os"

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
	fmt.Scanln(&prompt)

	for {

		if prompt == "" {
			fmt.Println("No prompt provided")
			break
		} else {
			response = chat(userID, prompt)
			color.Green("Model: %s", response)
			prompt = ""
			fmt.Printf("You: ")
			prompt, _ = reader.ReadString('\n')

		}

	}

}
