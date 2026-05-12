package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
)

func DanhGiaRoutes(r *gin.Engine,hub *websocket.Hub) {
	ctrl := controllers.NewDanhGiaController(hub)

	danhGia := r.Group("/danh-gia")
	{
		danhGia.POST("", ctrl.CreateDanhGia)
		danhGia.GET("", ctrl.GetDanhSachDanhGia)
		danhGia.GET("/:id", ctrl.GetDanhGiaByID)
		danhGia.PUT("/:id", ctrl.UpdateDanhGia)
		danhGia.DELETE("/:id", ctrl.DeleteDanhGia)
	}
}
