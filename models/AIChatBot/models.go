package core

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type VectorResult struct {
	ID       string         `json:"id"`
	Document string         `json:"document"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Distance *float64       `json:"distance,omitempty"`
}
