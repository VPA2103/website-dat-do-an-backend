package rag

import (
	"context"
	"strings"

	core "github.com/vpa/quanlynhahang-backend/models/AIChatBot"
)

type Service struct {
	llm    core.Gemini
	vector core.VectorStore
	fs     core.FileStore
}

func New(llm core.Gemini, vector core.VectorStore, fs core.FileStore) *Service {
	return &Service{llm: llm, vector: vector, fs: fs}
}
// func (s *Service) RetrieveContext(ctx context.Context, restaurantID string, query string, nResults int) (string, error) {
// 	query = strings.TrimSpace(query)
// 	if query == "" {
// 		return "", nil
// 	}

// 	restaurant, _ := s.fs.GetRestaurant(restaurantID)
// 	restaurantDoc := ingestRestaurantToDocument(restaurant)

// 	lines := []string{}
// 	if strings.TrimSpace(restaurantDoc) != "" {
// 		lines = append(lines, strings.TrimSpace(restaurantDoc))
// 	}

// 	emb, err := s.llm.Embed(ctx, query)
// 	if err != nil {
// 		// Best-effort: still return restaurant info if we have it.
// 		return strings.TrimSpace(strings.Join(lines, "\n")), nil
// 	}

// 	results := []core.VectorResult{}
// 	if len(emb) > 0 {
// 		results, err = s.vector.QueryMenu(ctx, restaurantID, emb, nResults)
// 		if err != nil {
// 			// Best-effort: still return restaurant info if we have it.
// 			return strings.TrimSpace(strings.Join(lines, "\n")), nil
// 		}
// 	}

// 	if len(results) > 0 {
// 		lines = append(lines, "\nMenu liên quan:")
// 		for _, r := range results {
// 			doc := strings.TrimSpace(r.Document)
// 			if doc != "" {
// 				lines = append(lines, "- "+doc)
// 			}
// 		}
// 	} else {
// 		// Fallback: if embeddings are missing/unavailable, provide a small menu list from DB.
// 		items, _, err := s.fs.ListMenuItems(restaurantID, 0, nResults)
// 		if err == nil && len(items) > 0 {
// 			lines = append(lines, "\nMenu (fallback, chưa có embedding):")
// 			for _, it := range items {
// 				doc := strings.TrimSpace(ingest.MenuItemToDocument(it))
// 				if doc != "" {
// 					lines = append(lines, "- "+doc)
// 				}
// 			}
// 		} else {
// 			lines = append(lines, "\nMenu liên quan: (không có dữ liệu menu trong hệ thống)")
// 		}
// 	}

// 	return strings.TrimSpace(strings.Join(lines, "\n")), nil
// }

func (s *Service) RetrieveContext(
	ctx context.Context,
	query string,
	nResults int,
) (string, error) {

	query = strings.TrimSpace(query)
	if query == "" {
		return "", nil
	}

	emb, err := s.llm.Embed(ctx, query)
	if err != nil {
		return "", nil
	}

	results, err := s.vector.QueryMenu(ctx, emb, nResults)
	if err != nil {
		return "", nil
	}

	if len(results) == 0 {
		return "", nil
	}

	var lines []string
	lines = append(lines, "MENU LIÊN QUAN:")

	for _, r := range results {
		lines = append(lines, "- "+r.Document)
	}

	return strings.Join(lines, "\n"), nil
}


