package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func BinhLuanRoutes(r *gin.Engine) {
	binhluan := r.Group("/binh-luan")
	{
		binhluan.POST("", middleware.AuthMiddleware(),controllers.CreateBinhLuan)
		binhluan.GET("/mon-an/:ma_mon_an", controllers.GetBinhLuanByMonAn)
		binhluan.GET("/:id",middleware.AuthMiddleware(), controllers.GetBinhLuanByID)
		binhluan.PUT("/:id", middleware.AuthMiddleware(),controllers.UpdateBinhLuan)
		binhluan.DELETE("/:id",middleware.AuthMiddleware(), controllers.DeleteBinhLuan)
	}
}