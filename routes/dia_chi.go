package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func DiaChiRoutes(r *gin.Engine) {
	dc := r.Group("/dia-chi")
	{
		dc.POST("/", middleware.AuthMiddleware(), controllers.CreateDiaChi)
		dc.GET("/user/:ma_nguoi_dung", middleware.AuthMiddleware(), controllers.GetDiaChiByUser)
		dc.GET("/:id", middleware.AuthMiddleware(), controllers.GetDiaChiByID)
		dc.PATCH("/:id", middleware.AuthMiddleware(), controllers.UpdateDiaChi)
		dc.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteDiaChi)
		dc.PATCH("/:id/mac-dinh", controllers.SetDiaChiMacDinh)
	}
}
