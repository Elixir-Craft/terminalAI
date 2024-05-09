package chat

import (
	"fmt"
	// uuid
	"github.com/google/uuid"
)

func ChatMode() {
	var prompt string
	var response string
	var userID string

	// uuid
	userID = uuid.New().String()
	fmt.Printf("You: ")
	fmt.Scanln(&prompt)

	for {

		if prompt == "" {
			fmt.Println("No prompt provided")
			break
		} else {
			response = chat(userID, prompt)
			fmt.Println("Model: ", response)
			prompt = ""
			fmt.Printf("You: ")
			fmt.Scanln(&prompt)
		}

	}

}
