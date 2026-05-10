package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/internal/repository"
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

	// 🧱 Tự động migrate
	err := config.DB.AutoMigrate(
		// core
		&models.NguoiDung{},
		&models.LoaiMonAn{},
		&models.MonAn{},

		// mid
		&models.BanAn{},
		&models.GiamGia{},
		&models.DiaChi{},

		// order/payment
		&models.DatBan{},
		&models.HoaDon{},
		&models.ChiTietHoaDon{},
		&models.ThanhToan{},

		// interaction (SAU CÙNG)
		&models.BinhLuan{},
		&models.DanhGia{},
		&models.YeuThich{},

		// misc
		&models.HinhAnh{},
		&models.ThongBao{},
		&models.LienHe{},
		&models.Message{},
		&models.Room{},
		&models.Payments{},
	)

	if err != nil {
		log.Fatalf("❌ Lỗi khi migrate DB: %v", err)
	}

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

	routes.SetupRoutes(r, chatUC, notiUC)

	handler := &websocket.Handler{
		ChatUC: chatUC,
		NotiUC: notiUC,
	}

	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleWS(hub, handler)(c.Writer, c.Request)
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
