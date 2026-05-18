package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func SePayPayment(r *gin.Engine) {

	r.POST("/payment/create", controllers.CreateSePayPaymentForm)

	// webhook
	r.POST("/hooks/sepay-payment", controllers.SePayWebhookHandler)

}
