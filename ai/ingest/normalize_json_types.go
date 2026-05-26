package ingest

import "github.com/vpa/quanlynhahang-backend/ai/core"

// Adapter types to avoid importing handler-private request structs.

type MenuItemJSONIn struct {
	ID          *string
	Name        string
	Description *string
	Price       *float64
	Tags        any
	Allergens   any
	Ingredients any
}

func NormalizeMenuItemFromJSONIn(in MenuItemJSONIn) core.MenuItem {
	return NormalizeMenuItemFromJSON(menuItemJSON{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Tags:        in.Tags,
		Allergens:   in.Allergens,
		Ingredients: in.Ingredients,
	})
}
