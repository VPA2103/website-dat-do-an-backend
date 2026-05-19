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
		ship.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("shipper"), ctrl.GetHoaDons)

		ship.PUT("/:id/trang-thai", middleware.AuthMiddleware(), middleware.RoleMiddleware("shipper"), ctrl.UpdateTrangThaiHoaDon)

	}
}
