package core

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MenuItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Allergens   []string `json:"allergens,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
}

type RestaurantInfo struct {
	Name      *string `json:"name,omitempty"`
	Address   *string `json:"address,omitempty"`
	OpenHours *string `json:"open_hours,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Style     *string `json:"style,omitempty"`
	Policies  *string `json:"policies,omitempty"`
}

type VectorResult struct {
	ID       string         `json:"id"`
	Document string         `json:"document"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Distance *float64       `json:"distance,omitempty"`
}