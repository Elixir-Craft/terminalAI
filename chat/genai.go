package chat

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	imageTemperature = 0.8
	chatTemperature  = 0.3
)

var geminiKey string

func createGeminiClient() (*genai.Client, error) {
	geminiKey = os.Getenv("GOOGLE_AI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}
	return client, nil
}

func startNewChatSession() *genai.ChatSession {
	ctx := context.Background()
	cs := createChatSession(ctx, os.Getenv("TERMINAL_AI_MODEL"), chatTemperature)
	return cs
}

func createChatSession(ctx context.Context, modelID string, temperature float32) *genai.ChatSession {
	client, err := createGeminiClient()
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel(modelID)
	model.Temperature = &temperature
	cs := model.StartChat()
	return cs
}

func send(cs *genai.ChatSession, msg string) *genai.GenerateContentResponse {
	if cs == nil {
		cs = startNewChatSession()
	}

	ctx := context.Background()
	res, err := cs.SendMessage(ctx, genai.Text(msg))

	if err != nil {
		log.Fatal(err)
	}
	return res
}

func printResponse(resp *genai.GenerateContentResponse) string {
	var ret string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			ret = ret + fmt.Sprintf("%v", part)
		}
	}
	return ret
}

// #################################

var userSessions = make(map[string]*genai.ChatSession)

func chat(userID string, msg string) string {
	var cs *genai.ChatSession
	if _, ok := userSessions[userID]; ok {
		cs = userSessions[userID]
	} else {
		cs = startNewChatSession()
		userSessions[userID] = cs
	}

	resp := send(cs, msg)
	return printResponse(resp)
}
