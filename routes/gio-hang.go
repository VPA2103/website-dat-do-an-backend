package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func GioHangRoutes(r *gin.Engine) {
	gioHang := r.Group("/gio-hang")
	{
		gioHang.POST("", middleware.AuthMiddleware(), controllers.AddToCart)

		gioHang.GET("", controllers.GetAllCart)

		gioHang.GET("/user/:id", controllers.GetCartByUser)

		gioHang.PUT("/:ma_mon_an", middleware.AuthMiddleware(), controllers.UpdateSoLuongCart)

		gioHang.DELETE("/:ma_mon_an", middleware.AuthMiddleware(), controllers.DeleteCart)

		gioHang.DELETE("/clear", middleware.AuthMiddleware(), controllers.XoaGioHangNguoiDung)
	}
}
