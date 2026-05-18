package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	auth := r.Group("/auth")
	{

		auth.POST("/send-otp", controllers.SendOTP)
		auth.POST("/reset-password", controllers.ResetPassword)
	}
}
