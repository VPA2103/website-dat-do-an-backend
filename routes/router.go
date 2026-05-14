package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/usecase"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func SetupRoutes(r *gin.Engine, chatUC *usecase.ChatUseCase,
	notiUC *usecase.NotificationUseCase, hub *websocket.Hub) {
	// 🌐 Route gốc
	r.GET("/", func(c *gin.Context) {
		c.File("./static/web-page.html")
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 🔒 Nhóm yêu cầu xác thực
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/profile", controllers.GetProfile)

	// 👑 Nhóm chỉ cho admin
	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/dashboard", controllers.AdminDashboard)
	contactHandler := &controllers.ContactHandler{
		NotiUC: notiUC,
	}
	// 👨‍💼 Nhân viên routes (có thể để ngoài hoặc trong nhóm admin)
	NguoiDungRoutes(r)
	BanAnRoutes(r)
	LoaiMonAnRoutes(r)
	MonAnRoutes(r)
	LienHeRoutes(r, contactHandler)
	DatBanRoutes(r)
	HoaDonRoutes(r)
	DiaChiRoutes(r)
	GiamGiaRoutes(r)
	BinhLuanRoutes(r, hub)
	YeuThichRoutes(r)
	DanhGiaRoutes(r, hub)
	GioHangRoutes(r)
	Payment(r)
	SePayPayment(r)
}
