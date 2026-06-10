package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func ShipRoutes(r *gin.Engine, hub *websocket.Hub) {

	ctrl := controllers.NewHoaDonController(hub)

	ship := r.Group("/ship")
	{
		ship.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("shipper"), ctrl.GetHoaDonByShipper)
		ship.GET("/all-hoa-don", middleware.AuthMiddleware(), middleware.RoleMiddleware("shipper"), ctrl.GetALLHoaDonByShipper)

		// ship.PUT("/:id/trang-thai", middleware.AuthMiddleware(), middleware.RoleMiddleware("shipper"), ctrl.UpdateTrangThaiHoaDon)

		// ship.POST("/accept-order", middleware.AuthMiddleware(), middleware.RoleMiddleware("shipper"), controllers.AcceptShipOrder)
	}
}
