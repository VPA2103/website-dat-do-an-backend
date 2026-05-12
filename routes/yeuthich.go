package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func YeuThichRoutes(r *gin.Engine) {
	yeuThich := r.Group("/yeu-thich")
	{
		yeuThich.POST("", middleware.AuthMiddleware(), controllers.AddMonAnYeuThich)
		yeuThich.GET("", controllers.GetAllYeuThich)
		yeuThich.GET("/user/:id", controllers.GetYeuThichByUser)
		yeuThich.DELETE("/:ma_mon_an", middleware.AuthMiddleware(), controllers.DeleteYeuThich)
	}
}
