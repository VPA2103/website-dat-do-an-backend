package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		panic("❌ DB_URL không tồn tại! Hãy kiểm tra Variables trong Railway.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect to database: %v", err))
	}

	DB = db
	fmt.Println("🚀 Database connected successfully")
}

func SetupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		AllowWebSockets: true,
		MaxAge:           12 * time.Hour,
	}))
}

// func ConnectDB() {
// 	host := os.Getenv("POSTGRES_HOST")
// 	port := os.Getenv("POSTGRES_PORT")
// 	user := os.Getenv("POSTGRES_USER")
// 	password := os.Getenv("POSTGRES_PASSWORD")
// 	dbname := os.Getenv("POSTGRES_DBNAME")
// 	sslmode := os.Getenv("POSTGRES_SSLMODE")

// 	// kiểm tra thiếu env
// 	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
// 		panic("❌ Thiếu biến môi trường Postgres")
// 	}

// 	dsn := fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
// 		host, port, user, password, dbname, sslmode,
// 	)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(fmt.Sprintf("❌ Failed to connect to database: %v", err))
// 	}

// 	DB = db
// 	fmt.Println("🚀 Database connected successfully")
// }
