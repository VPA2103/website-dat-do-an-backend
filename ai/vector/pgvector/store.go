package pgvector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	v "github.com/pgvector/pgvector-go"
	"github.com/vpa/quanlynhahang-backend/ai/config"
	"github.com/vpa/quanlynhahang-backend/ai/core"
	"github.com/vpa/quanlynhahang-backend/ai/vector/metric"
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

func (s *Store) UpsertMenuItem(ctx context.Context, restaurantID string, item core.MenuItem, embedding []float32, document string, metadata map[string]any) error {
	if err := s.validateEmbedding(embedding); err != nil {
		return err
	}
	if item.ID == "" {
		return errors.New("missing menu item id")
	}
	if document == "" {
		return errors.New("missing document")
	}
	b, _ := json.Marshal(metadata)
	cmd, err := s.pool.Exec(ctx, `
UPDATE menu_items
SET embedding = $2, document = $3, metadata = $4, updated_at = now()
WHERE id = $1
`, item.ID, v.NewVector(embedding), document, b)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("menu item not found")
	}
	return err
}

func (s *Store) DeleteMenuItem(ctx context.Context, restaurantID, id string) error {
	_, err := s.pool.Exec(ctx, `UPDATE menu_items SET embedding = NULL, document = NULL, metadata = NULL, updated_at = now() WHERE id = $1`, id)
	return err
}

func (s *Store) UpsertRestaurant(ctx context.Context, restaurantID string, embedding []float32, document string, metadata map[string]any) error {
	if err := s.validateEmbedding(embedding); err != nil {
		return err
	}
	if restaurantID == "" {
		return errors.New("missing restaurant id")
	}
	if document == "" {
		return errors.New("missing document")
	}
	b, _ := json.Marshal(metadata)
	cmd, err := s.pool.Exec(ctx, `
UPDATE restaurants
SET embedding = $2, document = $3, metadata = $4, updated_at = now()
WHERE id = $1
`, restaurantID, v.NewVector(embedding), document, b)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		// Create row if it does not exist.
		_, err = s.pool.Exec(ctx, `
INSERT INTO restaurants (id, embedding, document, metadata, updated_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (id) DO UPDATE SET embedding = EXCLUDED.embedding, document = EXCLUDED.document, metadata = EXCLUDED.metadata, updated_at = now()
`, restaurantID, v.NewVector(embedding), document, b)
	}
	return err
}

func (s *Store) QueryMenu(ctx context.Context, restaurantID string, embedding []float32, nResults int) ([]core.VectorResult, error) {
	if err := s.validateEmbedding(embedding); err != nil {
		return nil, err
	}
	if nResults <= 0 {
		nResults = 6
	}
	_, distOp, err := metric.Spec(s.vec.Metric)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`
SELECT id, document, metadata, (embedding %s $1) AS distance
FROM menu_items
WHERE embedding IS NOT NULL
ORDER BY embedding %s $1
LIMIT $2
`, distOp, distOp)

	rows, err := s.pool.Query(ctx, q, v.NewVector(embedding), nResults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []core.VectorResult{}
	for rows.Next() {
		var id, doc string
		var metaBytes []byte
		var dist float64
		if err := rows.Scan(&id, &doc, &metaBytes, &dist); err != nil {
			return nil, err
		}
		meta := map[string]any{}
		if len(metaBytes) > 0 {
			_ = json.Unmarshal(metaBytes, &meta)
		}
		d := dist
		out = append(out, core.VectorResult{ID: id, Document: doc, Metadata: meta, Distance: &d})
	}
	return out, rows.Err()
}
