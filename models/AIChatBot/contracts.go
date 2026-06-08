package core

import "context"

type FileStore interface {
	EnsureThread(threadID string) (string, error)
	AppendThreadMessage(threadID, role, content string) error
	GetThreadMessages(threadID string) ([]Message, error)
	//DeleteMenuItem(id string) (bool, error)
}

type VectorStore interface {
	//UpsertMenuItem(ctx context.Context, restaurantID string, item MenuItem, embedding []float32, document string, metadata map[string]any) error
	//DeleteMenuItem(ctx context.Context,id string) error
	//UpsertRestaurant(ctx context.Context,string, embedding []float32, document string, metadata map[string]any) error
	QueryMenu(ctx context.Context, embedding []float32, nResults int) ([]VectorResult, error)
}

type Gemini interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Generate(ctx context.Context, contents []string) (string, *RateLimitError, error)
}

type RAG interface {
	RetrieveContext(ctx context.Context, query string, nResults int) (string, error)
}

type RateLimitError struct {
	Message           string `json:"message"`
	RetryAfterSeconds *int   `json:"retry_after_seconds"`
}
