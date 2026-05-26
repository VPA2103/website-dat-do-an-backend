package ingest

import (
	"strconv"
	"strings"

	"github.com/vpa/quanlynhahang-backend/ai/core"
)

func NormalizeMenuItemFromMap(raw map[string]string) core.MenuItem {
	out := core.MenuItem{}
	if v := strings.TrimSpace(raw["id"]); v != "" {
		out.ID = v
	}
	out.Name = strings.TrimSpace(raw["name"])
	if v := strings.TrimSpace(raw["description"]); v != "" {
		out.Description = &v
	}
	if v := strings.TrimSpace(raw["price"]); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			out.Price = &f
		}
	}
	out.Tags = toList(raw["tags"])
	out.Allergens = toList(raw["allergens"])
	out.Ingredients = toList(raw["ingredients"])
	return out
}
