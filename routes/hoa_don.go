package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func HoaDonRoutes(r *gin.Engine) {

	hoaDon := r.Group("/hoa-don")
	{
		hoaDon.POST("", middleware.AuthMiddleware(), controllers.DatDoAn)

		hoaDon.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetHoaDons)

		hoaDon.GET("/:id", controllers.GetHoaDonByID)

		hoaDon.PUT("/:id", controllers.UpdateHoaDon)

		hoaDon.PUT("/:id/trang-thai", controllers.UpdateTrangThaiHoaDon)

		hoaDon.GET("/user", middleware.AuthMiddleware(), controllers.GetHoaDonByNguoiDung)

		hoaDon.PUT("/:id/huy", controllers.HuyHoaDon)

		hoaDon.GET("/trang-thai", controllers.GetHoaDonByTrangThai)

		hoaDon.DELETE("/:id", controllers.XoaHoaDon)

		// hoaDon.POST("/:ma_hd/thanh-toan", controllers.ThanhToanHoaDon)
	}
}
