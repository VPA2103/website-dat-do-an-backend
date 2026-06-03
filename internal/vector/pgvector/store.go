package pgvector

import (
	"fmt"

	"github.com/vpa/quanlynhahang-backend/config"

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
