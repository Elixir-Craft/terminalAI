package models

import (
	"context"
	"io"
	"log"
)

type Model StartChat
type StartChat func() Chat
type Chat func(ctx context.Context, prompt string) (StreamingOutput, error)

type StreamingOutput <-chan string

func (stream StreamingOutput) WriteTo(w io.Writer) (int64, error) {
	var nr int64
	for chunk := range stream {
		n, err := w.Write([]byte(chunk))
		nr += int64(n)
		if err != nil {
			return nr, err
		}
	}
	return nr, nil
}

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
