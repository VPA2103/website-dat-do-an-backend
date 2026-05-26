package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vpa/quanlynhahang-backend/ai/config"
	"github.com/vpa/quanlynhahang-backend/ai/vector/metric"
)

// EnsureSchema creates/updates the required DB schema for backend.golang.
// This is intentionally idempotent (CREATE IF NOT EXISTS).
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool, embDim int, vecCfg config.VectorConfig) error {
	if embDim <= 0 {
		return fmt.Errorf("invalid embedding dim: %d", embDim)
	}

	// pgvector extension
	if _, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector;`); err != nil {
		return err
	}

	restaurantsSQL := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS restaurants (
  id TEXT PRIMARY KEY,
  name TEXT,
  address TEXT,
  open_hours TEXT,
  phone TEXT,
  style TEXT,
  policies TEXT,
  embedding vector(%d),
  document TEXT,
  metadata JSONB,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`, embDim)
	if _, err := pool.Exec(ctx, restaurantsSQL); err != nil {
		return err
	}
	if err := ensureEmbeddingDim(ctx, pool, "restaurants", "embedding", embDim); err != nil {
		return err
	}

	menuSQL := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS menu_items (
  id TEXT PRIMARY KEY,
  restaurant_id TEXT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  price NUMERIC,
  tags JSONB,
  allergens JSONB,
  ingredients JSONB,
  embedding vector(%d),
  document TEXT,
  metadata JSONB,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`, embDim)
	if _, err := pool.Exec(ctx, menuSQL); err != nil {
		return err
	}
	if err := ensureEmbeddingDim(ctx, pool, "menu_items", "embedding", embDim); err != nil {
		return err
	}

	threadsSQL := `
CREATE TABLE IF NOT EXISTS threads (
  id TEXT PRIMARY KEY,
  restaurant_id TEXT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`
	if _, err := pool.Exec(ctx, threadsSQL); err != nil {
		return err
	}

	msgsSQL := `
CREATE TABLE IF NOT EXISTS thread_messages (
  id BIGSERIAL PRIMARY KEY,
  thread_id TEXT NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
  restaurant_id TEXT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`
	if _, err := pool.Exec(ctx, msgsSQL); err != nil {
		return err
	}

	// Useful btree indexes.
	if _, err := pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_menu_items_restaurant_id ON menu_items (restaurant_id);`); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_thread_messages_thread_created ON thread_messages (thread_id, created_at);`); err != nil {
		return err
	}

	// Vector indexes (HNSW) for faster similarity search.
	// NOTE: HNSW indexing is only helpful when there are enough rows; for small datasets it is fine either way.
	if vecCfg.Index != "none" {
		opClass, _, err := metric.Spec(vecCfg.Metric)
		if err != nil {
			return err
		}
		if vecCfg.Index == "hnsw" {
			idxMenu := fmt.Sprintf("idx_menu_items_embedding_hnsw_%s", vecCfg.Metric)
			if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON menu_items USING hnsw (embedding %s);`, idxMenu, opClass)); err != nil {
				return err
			}
			idxRest := fmt.Sprintf("idx_restaurants_embedding_hnsw_%s", vecCfg.Metric)
			if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON restaurants USING hnsw (embedding %s);`, idxRest, opClass)); err != nil {
				return err
			}
		} else if vecCfg.Index == "ivfflat" {
			idxMenu := fmt.Sprintf("idx_menu_items_embedding_ivfflat_%s", vecCfg.Metric)
			if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON menu_items USING ivfflat (embedding %s) WITH (lists = %d);`, idxMenu, opClass, vecCfg.IVFFlatLists)); err != nil {
				return err
			}
			idxRest := fmt.Sprintf("idx_restaurants_embedding_ivfflat_%s", vecCfg.Metric)
			if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON restaurants USING ivfflat (embedding %s) WITH (lists = %d);`, idxRest, opClass, vecCfg.IVFFlatLists)); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported VECTOR_INDEX: %q", vecCfg.Index)
		}
	}

	return nil
}

// ensureEmbeddingDim keeps the vector dimension of an embedding column aligned with the configured embDim.
// Embeddings are derived data; if the dimension changes, we drop existing values (set to NULL) and alter the column.
func ensureEmbeddingDim(ctx context.Context, pool *pgxpool.Pool, table, column string, embDim int) error {
	// Only allow known identifiers (avoid SQL injection).
	if (table != "restaurants" && table != "menu_items") || column != "embedding" {
		return fmt.Errorf("unsupported embedding column: %s.%s", table, column)
	}

	var typmod int
	err := pool.QueryRow(ctx, `
SELECT atttypmod
FROM pg_attribute
WHERE attrelid = $1::regclass AND attname = $2
`, table, column).Scan(&typmod)
	if err != nil {
		return err
	}
	if typmod == embDim {
		return nil
	}

	// Drop any known vector indexes that might depend on the column.
	metrics := []string{"cosine", "l2", "ip"}
	indexes := []string{"hnsw", "ivfflat"}
	for _, idx := range indexes {
		for _, metric := range metrics {
			name := fmt.Sprintf("idx_%s_embedding_%s_%s", table, idx, metric)
			_, _ = pool.Exec(ctx, fmt.Sprintf(`DROP INDEX IF EXISTS %s;`, name))
		}
	}

	// Embeddings are derived; clear existing values then change dim.
	_, _ = pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET %s = NULL WHERE %s IS NOT NULL;`, table, column, column))
	_, err = pool.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s TYPE vector(%d);`, table, column, embDim))
	return err
}
