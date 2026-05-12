package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func HoaDonRoutes(r *gin.Engine) {

	hoaDon := r.Group("/hoa-don")
	{
		hoaDon.POST("", controllers.DatDoAn)

		hoaDon.GET("", controllers.GetHoaDons)

		hoaDon.GET("/:id", controllers.GetHoaDonByID)

		hoaDon.PUT("/:id", controllers.UpdateHoaDon)

		hoaDon.PUT("/:id/trang-thai", controllers.UpdateTrangThaiHoaDon)

		hoaDon.PUT("/:id/huy", controllers.HuyHoaDon)

		hoaDon.GET("/trang-thai", controllers.GetHoaDonByTrangThai)

		hoaDon.DELETE("/:id", controllers.XoaHoaDon)

		hoaDon.POST("/:ma_hd/thanh-toan", controllers.ThanhToanHoaDon)
	}
}
