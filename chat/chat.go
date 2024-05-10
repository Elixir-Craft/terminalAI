package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"terminalAI/models"

	"github.com/fatih/color"
)

func ChatMode(model models.Model) {
	reader := bufio.NewReader(os.Stdin)
	chat := model()

	for {
		fmt.Printf("You: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)

		if prompt == "/exit" {
			color.Cyan("Goodbye!")
			break
		} else if prompt == "" {
			color.Red("No prompt provided")
		} else {
			response, err := chat(context.Background(), prompt)
			if err != nil {
				log.Fatal(err)
			}

			func() {
				color.Set(color.FgGreen)
				defer color.Unset()
				fmt.Print("Model: ")
				_, err = response.WriteTo(os.Stdout)
			}()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
