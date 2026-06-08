package pgvector

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/controllers"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
	dim  int
	vec  config.VectorConfig
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// NewStoreWithConfig constructs a Store with the configured embedding dimension and vector settings.
// Schema and index creation is owned by internal/db.EnsureSchema.
func NewStoreWithConfig(pool *pgxpool.Pool, dim int, vecCfg config.VectorConfig) *Store {
	return &Store{pool: pool, dim: dim, vec: vecCfg}
}

func (s *Store) validateEmbedding(embedding []float32) error {
	if s.dim <= 0 {
		return nil
	}
	if len(embedding) != s.dim {
		return fmt.Errorf("embedding dim mismatch: got=%d want=%d", len(embedding), s.dim)
	}
	return nil
}
func (s *Store) QueryMenu(
	ctx context.Context,
	embedding []float32,
	nResults int,
) ([]controllers.VectorResult, error) {

	vec := vectorToString(embedding)

	query := `
SELECT id, document, metadata,
       embedding <=> $1 AS distance
FROM menu_embeddings
ORDER BY embedding <=> $1
LIMIT $2
`

	rows, err := s.pool.Query(ctx, query, vec, nResults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []controllers.VectorResult{}

	for rows.Next() {
		var r controllers.VectorResult
		var meta []byte
		var dist float64

		if err := rows.Scan(&r.ID, &r.Document, &meta, &dist); err != nil {
			return nil, err
		}

		r.Distance = &dist
		_ = json.Unmarshal(meta, &r.Metadata)

		results = append(results, r)
	}

	return results, nil
}

func vectorToString(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}

	var b strings.Builder
	b.WriteString("[")

	for i, v := range vec {
		b.WriteString(fmt.Sprintf("%f", v))
		if i < len(vec)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")
	return b.String()
}
