package core

import "context"

type FileStore interface {
	EnsureThread(restaurantID, threadID string) (string, error)
	AppendThreadMessage(restaurantID, threadID, role, content string) error
	GetThreadMessages(restaurantID, threadID string) ([]Message, error)

	UpsertMenuItem(restaurantID string, item MenuItem) (MenuItem, error)
	ListMenuItems(restaurantID string, offset, limit int) ([]MenuItem, int, error)
	DeleteMenuItem(restaurantID, id string) (bool, error)

	SetRestaurant(restaurantID string, info RestaurantInfo) (RestaurantInfo, error)
	GetRestaurant(restaurantID string) (RestaurantInfo, error)
}

type VectorStore interface {
	UpsertMenuItem(ctx context.Context, restaurantID string, item MenuItem, embedding []float32, document string, metadata map[string]any) error
	DeleteMenuItem(ctx context.Context, restaurantID, id string) error
	UpsertRestaurant(ctx context.Context, restaurantID string, embedding []float32, document string, metadata map[string]any) error
	QueryMenu(ctx context.Context, restaurantID string, embedding []float32, nResults int) ([]VectorResult, error)
}

type Gemini interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Generate(ctx context.Context, contents []string) (string, *RateLimitError, error)
}

type RAG interface {
	RetrieveContext(ctx context.Context, restaurantID string, query string, nResults int) (string, error)
}

type RateLimitError struct {
	Message           string `json:"message"`
	RetryAfterSeconds *int   `json:"retry_after_seconds"`
}