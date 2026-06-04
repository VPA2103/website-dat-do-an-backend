package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func NhaHangRoutes(r *gin.Engine,h*controllers.ChatHandler) {
	nhahang := r.Group("/nha-hang")
	{
		nhahang.POST("/create",middleware.AuthMiddleware(), h.CreateNhaHang)

		nhahang.GET("/all", controllers.GetAllNhaHang)

		nhahang.GET("/:id", controllers.GetNhaHangByID)

		nhahang.GET("/user",middleware.AuthMiddleware(), controllers.GetNhaHangByUser)

		nhahang.PATCH("/update/:id",middleware.AuthMiddleware(),controllers.UpdateNhaHang)

		nhahang.DELETE("/delete/:id",middleware.AuthMiddleware(), controllers.DeleteNhaHang)
	}
}