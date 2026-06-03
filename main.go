package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/llm"
	"github.com/vpa/quanlynhahang-backend/internal/repository"
	"github.com/vpa/quanlynhahang-backend/internal/store"
	"github.com/vpa/quanlynhahang-backend/internal/usecase"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/routes"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠Không tìm thấy file .env, dùng SECRET_KEY mặc định")
	}
	config.LoadPaymentConfig()
	// 💾 Kết nối Cloudinary
	config.InitCloudinary()
	// 🔧 Khởi tạo Gin
	r := gin.Default()

	// ⚙️ Cấu hình CORS
	config.SetupCORS(r)

	// 💾 Kết nối DB
	config.ConnectDB()

	config.DB.DisableForeignKeyConstraintWhenMigrating = true
	// Bước 1: Migrate bảng cha riêng, kiểm tra lỗi
	err := config.DB.AutoMigrate(
		&models.NguoiDung{},
		&models.LoaiMonAn{},
		&models.MonAn{},
		&models.NhomOption{},
		&models.OptionItem{},
	)
	if err != nil {
		log.Fatal("❌1 Migrate core failed:", err)
	}

	// Bước 2: Migrate bảng con sau khi chắc chắn bảng cha đã tồn tại
	err = config.DB.AutoMigrate(
		&models.BanAn{},
		&models.GiamGia{},
		&models.DiaChi{},
		&models.DatBan{},
		&models.HoaDon{},
		&models.ChiTietHoaDon{},
		&models.ChiTietHoaDonOption{},
		&models.ThanhToan{},
		&models.GioHang{},

		&models.BinhLuan{},
		&models.DanhGia{},
		&models.YeuThich{},

		&models.HinhAnh{},
		&models.ThongBao{},
		&models.LienHe{},
		&models.Message{},
		&models.Room{},
		&models.Payments{},
		&models.GioHangOption{},
		&models.NhaHang{},
	)
	if err != nil {
		log.Fatal("❌2 Migrate relations failed:", err)
	}
	log.Println("✅ Database migrated")

	geminiCfg := config.LoadGeminiConfig()

	// 2️⃣ Init Gemini LLM
	
	geminiLLM, err := llm.NewGemini(geminiCfg)
	if err != nil {
		log.Fatal("❌ Gemini init failed:", err)
	}
	

	// 3️⃣ PGX pool (dùng cho thread + message)
	pgxURL := os.Getenv("PGX_DATABASE_URL")
	if pgxURL == "" {
		log.Fatal("❌ Missing PGX_DATABASE_URL")
	}

	pgxPool, err := pgxpool.New(context.Background(), pgxURL)
	if err != nil {
		log.Fatal("❌ PGX pool init failed:", err)
	}

	// 4️⃣ FileStore
	fileStore := store.NewPostgresStore(pgxPool)

	// 5️⃣ Chat handler (❗ RAG = nil)
	chatHandler := controllers.NewChatHandler(
		fileStore,
		nil, // 👈 KHÔNG dùng RAG
		geminiLLM,
	)

	// 8️⃣ Register chatbot route
	routes.RegisterRoutes(r, chatHandler)

	// 🚏 Đăng ký route
	routes.UploadRoutes(r)

	//realtime
	hub := websocket.NewHub()
	go hub.Run()

	repo := repository.NewMessageRepository(config.DB)
	notiRepo := repository.NewNotificationRepository(config.DB)

	chatUC := &usecase.ChatUseCase{
		RT:   hub,
		Repo: repo,
	}
	notiUC := &usecase.NotificationUseCase{
		RT:   hub,
		Repo: notiRepo,
	}

	routes.SetupRoutes(r, chatUC, notiUC, hub)

	handler := &websocket.Handler{
		ChatUC: chatUC,
		NotiUC: notiUC,
	}

	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleWS(hub, handler)(c.Writer, c.Request)
	})

	r.GET("/ws/public", func(c *gin.Context) {
		websocket.HandleWSPublic(hub)(c.Writer, c.Request)
	})

	// 🚀 Chạy server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // chạy local
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Không thể khởi chạy server: %v", err)
	}

}
