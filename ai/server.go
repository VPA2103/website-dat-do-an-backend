package ai

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/ai/config"
	"github.com/vpa/quanlynhahang-backend/ai/db"
	"github.com/vpa/quanlynhahang-backend/ai/httpserver"
	"github.com/vpa/quanlynhahang-backend/ai/llm"
	"github.com/vpa/quanlynhahang-backend/ai/rag"
	"github.com/vpa/quanlynhahang-backend/ai/store"
	"github.com/vpa/quanlynhahang-backend/ai/vector/pgvector"
)

// RegisterRoutes khởi tạo toàn bộ chatbot và mount vào Gin group /ai
func RegisterRoutes(r *gin.Engine) {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("chatbot config: %v", err)
    }

    ctx := context.Background()
    pool, err := db.NewPool(ctx, cfg.Postgres)
    if err != nil {
        log.Fatalf("chatbot postgres: %v", err)
    }

    if err := db.EnsureSchema(ctx, pool, cfg.EmbeddingDim, cfg.Vector); err != nil {
        log.Fatalf("chatbot schema: %v", err)
    }

    pg := pgvector.NewStoreWithConfig(pool, cfg.EmbeddingDim, cfg.Vector)
    fileStore := store.NewPostgresStore(pool)

    gemini, err := llm.NewGemini(cfg.Gemini)
    if err != nil {
        log.Fatalf("chatbot gemini: %v", err)
    }

    svc := &httpserver.Services{
        FileStore: fileStore,
        Vector:    pg,
        Gemini:    gemini,
        RAG:       rag.New(gemini, pg, fileStore),
    }

    // Lấy http.Handler gốc của chatbot
    chatHandler := httpserver.NewRouter(cfg, svc)

    // Mount vào Gin dưới prefix /ai
    // Bridge net/http -> Gin bằng gin.WrapH
    r.Any("/ai/*path", gin.WrapH(http.StripPrefix("/ai", chatHandler)))
}