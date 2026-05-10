package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func BinhLuanRoutes(r *gin.Engine) {
	binhluan := r.Group("/binh-luan")
	{
		binhluan.POST("", controllers.CreateBinhLuan)
		binhluan.GET("/mon-an/:ma_mon_an", controllers.GetBinhLuanByMonAn)
		binhluan.GET("/:id", controllers.GetBinhLuanByID)
		binhluan.PUT("/:id", controllers.UpdateBinhLuan)
		binhluan.DELETE("/:id", controllers.DeleteBinhLuan)
	}
}