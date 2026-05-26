package ingest

import (
	"context"

	"github.com/vpa/quanlynhahang-backend/ai/core"
)

func MenuItemToDocument(item core.MenuItem) string {
	return menuItemToDocument(item)
}

func EmbedMenuItem(ctx context.Context, restaurantID string, item core.MenuItem, vector core.VectorStore, llm core.Gemini) error {
	doc := menuItemToDocument(item)
	emb, err := llm.Embed(ctx, doc)
	if err != nil {
		return err
	}
	if len(emb) == 0 {
		return nil
	}
	meta := map[string]any{"name": item.Name}
	if item.Price != nil {
		meta["price"] = *item.Price
	}
	if len(item.Tags) > 0 {
		meta["tags"] = item.Tags
	}
	return vector.UpsertMenuItem(ctx, restaurantID, item, emb, doc, meta)
}

func EmbedRestaurant(ctx context.Context, restaurantID string, info core.RestaurantInfo, vector core.VectorStore, llm core.Gemini) error {
	doc := restaurantToDocument(info)
	emb, err := llm.Embed(ctx, doc)
	if err != nil {
		return err
	}
	if len(emb) == 0 {
		return nil
	}
	meta := map[string]any{}
	if info.Name != nil {
		meta["name"] = *info.Name
	}
	return vector.UpsertRestaurant(ctx, restaurantID, emb, doc, meta)
}
