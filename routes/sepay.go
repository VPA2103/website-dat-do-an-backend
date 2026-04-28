package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func SePayPayment(r *gin.Engine) {
	r.POST("/api/payment/create", controllers.CreatePayment)
	r.POST("/api/payment/ipn", controllers.HandleIPN)

	// r.GET("/payment/success", SuccessPage)
	// r.GET("/payment/cancel", CancelPage)
	// r.GET("/payment/error", ErrorPage)
}
