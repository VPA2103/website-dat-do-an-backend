package routes

import (
	"github.com/vpa/quanlynhahang-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, chatbot *controllers.ChatHandler) {
	api := r.Group("/api")
	{
		api.POST("/chat", chatbot.Chat)
	}
}