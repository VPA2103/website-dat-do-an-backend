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

	config.DB.DisableForeignKeyConstraintWhenMigrating = true
	// Bước 1: Migrate bảng cha riêng, kiểm tra lỗi
	err := config.DB.AutoMigrate(
		&models.NguoiDung{},
		&models.LoaiMonAn{},
		&models.MonAn{},
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
	)
	if err != nil {
		log.Fatal("❌2 Migrate relations failed:", err)
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

	routes.SetupRoutes(r, chatUC, notiUC, hub)

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
