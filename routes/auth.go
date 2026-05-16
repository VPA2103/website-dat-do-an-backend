package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func AuthRoutes(r *gin.Engine) {

	auth := r.Group("/auth")
	{
		auth.POST("/send-otp", controllers.SendOTP)
		auth.POST("/reset-password", controllers.ResetPassword)
	}
}
