package routes

import (


	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/usecase"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func SetupRoutes(r *gin.Engine, chatUC *usecase.ChatUseCase, 
	notiUC *usecase.NotificationUseCase, hub *websocket.Hub,chatHandler *controllers.ChatHandler,) {
	// 🌐 Route gốc
	r.GET("/", func(c *gin.Context) {
		c.File("./static/web-page.html")
	})

	// 🔒 Nhóm yêu cầu xác thực
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/profile", controllers.GetProfile)

	// 👑 Nhóm chỉ cho admin
	// admin := auth.Group("/admin")
	// admin.Use(middleware.RoleMiddleware("admin"))
	// admin.GET("/dashboard", controllers.AdminDashboard)
	contactHandler := &controllers.ContactHandler{
		NotiUC: notiUC,
	}

	// 👨‍💼 Nhân viên routes (có thể để ngoài hoặc trong nhóm admin)
	AuthRoutes(r)

	NguoiDungRoutes(r,hub)

	BanAnRoutes(r)

	LoaiMonAnRoutes(r)

	MonAnRoutes(r, chatHandler)

	LienHeRoutes(r, contactHandler)

	DatBanRoutes(r)

	HoaDonRoutes(r, hub)

	DiaChiRoutes(r)

	GiamGiaRoutes(r)

	BinhLuanRoutes(r, hub)

	YeuThichRoutes(r)

	DanhGiaRoutes(r, hub)

	GioHangRoutes(r)

	Payment(r)

	SePayPayment(r, hub)

	OptionRoutes(r)

	ShipRoutes(r, hub)

	NhaHangRoutes(r,chatHandler)

	UserRoutes(r)

	AdminRoutes(r)
}
