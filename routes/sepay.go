package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
)

func SePayPayment(r *gin.Engine, hub *websocket.Hub) {

	ctrl := controllers.NewThanhToanController(hub)

	r.POST("/payment/create", controllers.CreateSePayPaymentForm)

	// webhook
	r.POST("/hooks/sepay-payment", ctrl.SePayWebhook)

}
