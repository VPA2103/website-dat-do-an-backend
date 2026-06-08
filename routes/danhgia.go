package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
)

func DanhGiaRoutes(r *gin.Engine, hub *websocket.Hub) {
	ctrl := controllers.NewDanhGiaController(hub)

	danhGia := r.Group("/danh-gia")
	{
		danhGia.POST("", ctrl.CreateDanhGia)
		danhGia.GET("/mon", controllers.GetRatingByMon)
		danhGia.GET("/mon/:id", controllers.GetDanhGiaByMonAn)
		danhGia.PUT("/:id", ctrl.UpdateDanhGia)
		danhGia.DELETE("/:id", ctrl.DeleteDanhGia)
		danhGia.GET("/check", controllers.CheckDanhGia)
		danhGia.GET("/so_luong_danh_gia", ctrl.GetSoLuongDanhGiaHomNay)
	}
}
