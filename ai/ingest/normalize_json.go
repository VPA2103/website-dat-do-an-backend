package ingest

import (
	"strings"

	"github.com/vpa/quanlynhahang-backend/ai/core"
)

type menuItemJSON struct {
	ID          *string
	Name        string
	Description *string
	Price       *float64
	Tags        any
	Allergens   any
	Ingredients any
}

func NormalizeMenuItemFromJSON(in menuItemJSON) core.MenuItem {
	item := core.MenuItem{}
	if in.ID != nil {
		item.ID = strings.TrimSpace(*in.ID)
	}
	item.Name = strings.TrimSpace(in.Name)
	item.Description = in.Description
	item.Price = in.Price
	item.Tags = normalizeAnyList(in.Tags)
	item.Allergens = normalizeAnyList(in.Allergens)
	item.Ingredients = normalizeAnyList(in.Ingredients)
	return item
}
