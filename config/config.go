package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type GeminiConfig struct {
	APIKey         string
	Model          string
	EmbeddingModel string
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
type Config struct {
	Vector VectorConfig
	Gemini GeminiConfig
}
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// func ConnectDB() {
// 	dsn := os.Getenv("DB_URL")

// 	if dsn == "" {
// 		panic("❌ DB_URL không tồn tại! Hãy kiểm tra Variables trong Railway.")
// 	}

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(fmt.Sprintf("❌ Failed to connect to database: %v", err))
// 	}

// 	DB = db
// 	fmt.Println("🚀 Database connected successfully")
// }

func SetupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))
}

func ConnectDB() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")
	sslmode := os.Getenv("POSTGRES_SSLMODE")

	// kiểm tra thiếu env
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		panic("❌ Thiếu biến môi trường Postgres")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		host, port, user, password, dbname, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect to database: %v", err))
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS vector`).Error; err != nil {
		panic(fmt.Sprintf("❌ Failed to create extension vector: %v", err))
	}

	DB = db
	fmt.Println("✔ Database connected successfully")
}

func LoadGeminiConfig() GeminiConfig {
	return GeminiConfig{
		APIKey:         os.Getenv("GEMINI_API_KEY"),
		Model:          os.Getenv("GEMINI_MODEL"),
		EmbeddingModel: os.Getenv("GEMINI_EMBEDDING_MODEL"),
	}
}
