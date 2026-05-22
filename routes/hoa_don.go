package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func HoaDonRoutes(r *gin.Engine, hub *websocket.Hub) {

	ctrl := controllers.NewHoaDonController(hub)

	hoaDon := r.Group("/hoa-don")
	{
		hoaDon.POST("", middleware.AuthMiddleware(), ctrl.DatDoAn)

		hoaDon.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin" , "shipper" ), ctrl.GetHoaDons)

		hoaDon.PUT("/:id", ctrl.UpdateHoaDon)

		hoaDon.PUT("/:id/trang-thai", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "shipper"), ctrl.UpdateTrangThaiHoaDon)

		hoaDon.GET("/user", middleware.AuthMiddleware(), ctrl.GetHoaDonByNguoiDung)

		hoaDon.GET("/trang-thai", ctrl.GetHoaDonByTrangThai)

		hoaDon.GET("/:id", middleware.AuthMiddleware(), ctrl.GetHoaDonByID)

		hoaDon.PUT("/:id/huy", ctrl.HuyHoaDon)

		hoaDon.DELETE("/:id", ctrl.XoaHoaDon)

		// hoaDon.POST("/:ma_hd/thanh-toan", controllers.ThanhToanHoaDon)
	}
}
