package chat

import (
	"fmt"
)

func ChatMode() {
	var prompt string
	var response string
	var userID string
	userID = "123479862d"

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
