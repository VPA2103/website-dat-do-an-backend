package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func DanhGiaRoutes(r *gin.Engine) {
	danhGia := r.Group("/danh-gia")
	{
		danhGia.POST("", controllers.CreateDanhGia)
		danhGia.GET("", controllers.GetDanhSachDanhGia)
		danhGia.GET("/:id", controllers.GetDanhGiaByID)
		danhGia.PUT("/:id", controllers.UpdateDanhGia)
		danhGia.DELETE("/:id", controllers.DeleteDanhGia)
	}
}