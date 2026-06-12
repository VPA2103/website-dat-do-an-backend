package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func DanhGiaRoutes(r *gin.Engine, hub *websocket.Hub) {
	ctrl := controllers.NewDanhGiaController(hub)

	danhGia := r.Group("/danh-gia")
	{
		danhGia.PUT("/an-danh-gia/:id", ctrl.AnDanhGia)
		danhGia.PUT("/hien-danh-gia/:id", ctrl.HienDanhGia)
		danhGia.POST("", ctrl.CreateDanhGia)
		danhGia.GET("/mon", controllers.GetRatingByMon)
		danhGia.GET("/tat-ca-danh-gia/:id", controllers.GetAllDanhGiaByMonAn)
		danhGia.GET("/mon/:id", controllers.GetDanhGiaByMonAn)
		danhGia.PUT("/:id", ctrl.UpdateDanhGia)
		danhGia.DELETE("/:id",middleware.AuthMiddleware(), ctrl.DeleteDanhGia)
		danhGia.GET("/check", controllers.CheckDanhGia)
		danhGia.GET("/so_luong_danh_gia", ctrl.GetSoLuongDanhGiaHomNay)
		danhGia.GET("/get-all-danh-gia-by-nguoi-dung",middleware.AuthMiddleware(), ctrl.GetAllDanhGiaByNguoiDung)
	}
}
