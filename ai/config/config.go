package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	CorsOrigins  []string
	EmbeddingDim int
	Vector       VectorConfig
	Gemini       GeminiConfig
	Postgres     PostgresConfig
}

type VectorConfig struct {
	// Metric controls the distance function used for vector search.
	// Supported: "cosine" (default), "l2", "ip".
	Metric string
	// Index controls which pgvector index to create.
	// Supported: "none" (default), "hnsw", "ivfflat".
	Index string
	// IVFFlatLists controls the ivfflat index lists parameter (only used when Index=ivfflat).
	IVFFlatLists int
}

type GeminiConfig struct {
	APIKey         string
	Model          string
	EmbeddingModel string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (Config, error) {
	// Best-effort load ../.env for local dev, without overriding existing env.
	_ = godotenv.Load(filepath.Join("..", ".env"))
	_ = godotenv.Load(".env")

	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if apiKey == "" {
		return Config{}, errors.New("missing GEMINI_API_KEY")
	}

	pgPort, _ := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	embDim, _ := strconv.Atoi(getEnv("EMBEDDING_DIM", "3072"))
	ivfLists, _ := strconv.Atoi(getEnv("VECTOR_IVFFLAT_LISTS", "100"))

	corsOrigins := splitCSV(getEnv("CORS_ORIGINS", "http://localhost:3000"))

	cfg := Config{
		Port:         getEnv("PORT", "8080"),
		CorsOrigins:  corsOrigins,
		EmbeddingDim: embDim,
		Vector: VectorConfig{
			Metric:       strings.ToLower(getEnv("VECTOR_METRIC", "cosine")),
			Index:        strings.ToLower(getEnv("VECTOR_INDEX", "none")),
			IVFFlatLists: ivfLists,
		},
		Gemini: GeminiConfig{
			APIKey:         apiKey,
			Model:          getEnv("GEMINI_MODEL", "gemini-2.5-flash-lite"),
			EmbeddingModel: getEnv("GEMINI_EMBEDDING_MODEL", "gemini-embedding-2"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     pgPort,
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DBNAME", "chatbot"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
	}

	if cfg.EmbeddingDim <= 0 {
		return Config{}, fmt.Errorf("invalid EMBEDDING_DIM: %d", cfg.EmbeddingDim)
	}
	if cfg.Postgres.Port <= 0 {
		return Config{}, fmt.Errorf("invalid POSTGRES_PORT: %d", cfg.Postgres.Port)
	}
	if cfg.Vector.Metric != "cosine" && cfg.Vector.Metric != "l2" && cfg.Vector.Metric != "ip" {
		return Config{}, fmt.Errorf("invalid VECTOR_METRIC: %q", cfg.Vector.Metric)
	}
	if cfg.Vector.Index != "none" && cfg.Vector.Index != "hnsw" && cfg.Vector.Index != "ivfflat" {
		return Config{}, fmt.Errorf("invalid VECTOR_INDEX: %q", cfg.Vector.Index)
	}
	if cfg.Vector.IVFFlatLists <= 0 {
		return Config{}, fmt.Errorf("invalid VECTOR_IVFFLAT_LISTS: %d", cfg.Vector.IVFFlatLists)
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
