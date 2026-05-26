package rag

import (
	"context"
	"strings"

	"github.com/vpa/quanlynhahang-backend/ai/core"
	"github.com/vpa/quanlynhahang-backend/ai/ingest"
)

type Service struct {
	llm    core.Gemini
	vector core.VectorStore
	fs     core.FileStore
}

func New(llm core.Gemini, vector core.VectorStore, fs core.FileStore) *Service {
	return &Service{llm: llm, vector: vector, fs: fs}
}

func (s *Service) RetrieveContext(ctx context.Context, restaurantID string, query string, nResults int) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", nil
	}

	restaurant, _ := s.fs.GetRestaurant(restaurantID)
	restaurantDoc := ingestRestaurantToDocument(restaurant)

	lines := []string{}
	if strings.TrimSpace(restaurantDoc) != "" {
		lines = append(lines, strings.TrimSpace(restaurantDoc))
	}

	emb, err := s.llm.Embed(ctx, query)
	if err != nil {
		// Best-effort: still return restaurant info if we have it.
		return strings.TrimSpace(strings.Join(lines, "\n")), nil
	}

	results := []core.VectorResult{}
	if len(emb) > 0 {
		results, err = s.vector.QueryMenu(ctx, restaurantID, emb, nResults)
		if err != nil {
			// Best-effort: still return restaurant info if we have it.
			return strings.TrimSpace(strings.Join(lines, "\n")), nil
		}
	}

	if len(results) > 0 {
		lines = append(lines, "\nMenu liên quan:")
		for _, r := range results {
			doc := strings.TrimSpace(r.Document)
			if doc != "" {
				lines = append(lines, "- "+doc)
			}
		}
	} else {
		// Fallback: if embeddings are missing/unavailable, provide a small menu list from DB.
		items, _, err := s.fs.ListMenuItems(restaurantID, 0, nResults)
		if err == nil && len(items) > 0 {
			lines = append(lines, "\nMenu (fallback, chưa có embedding):")
			for _, it := range items {
				doc := strings.TrimSpace(ingest.MenuItemToDocument(it))
				if doc != "" {
					lines = append(lines, "- "+doc)
				}
			}
		} else {
			lines = append(lines, "\nMenu liên quan: (không có dữ liệu menu trong hệ thống)")
		}
	}

	return strings.TrimSpace(strings.Join(lines, "\n")), nil
}

func ingestRestaurantToDocument(info core.RestaurantInfo) string {
	// Keep exact format from Python restaurant_to_document
	lines := []string{}
	if info.Name != nil && *info.Name != "" {
		lines = append(lines, "name: "+*info.Name)
	}
	if info.Address != nil && *info.Address != "" {
		lines = append(lines, "address: "+*info.Address)
	}
	if info.OpenHours != nil && *info.OpenHours != "" {
		lines = append(lines, "open_hours: "+*info.OpenHours)
	}
	if info.Phone != nil && *info.Phone != "" {
		lines = append(lines, "phone: "+*info.Phone)
	}
	if info.Style != nil && *info.Style != "" {
		lines = append(lines, "style: "+*info.Style)
	}
	if info.Policies != nil && *info.Policies != "" {
		lines = append(lines, "policies: "+*info.Policies)
	}
	return "Thông tin nhà hàng\n" + strings.Join(lines, "\n")
}
