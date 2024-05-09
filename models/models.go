package models

import (
	"context"
	"log"
)

type Model StartChat
type StartChat func() Chat
type Chat func(ctx context.Context, prompt string) (string, error)

type Backend func(modelName string) Model

var backends = map[string]Backend{}

func RegisterBackend(name string, backend Backend) {
	backends[name] = backend
}

func NewModel(backendName string, modelName string) Model {
	backend, k := backends[backendName]
	if !k {
		log.Fatalf("Backend '%s' not found.", backendName)
	}
	model := backend(modelName)
	return model
}
