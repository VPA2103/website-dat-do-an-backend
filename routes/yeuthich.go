package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func YeuThichRoutes(r *gin.Engine) {
	yeuThich := r.Group("/yeu-thich")
	{
		yeuThich.POST("", controllers.CreateYeuThich)
		yeuThich.GET("", controllers.GetAllYeuThich)
		yeuThich.GET("/user/:id", controllers.GetYeuThichByUser)
		yeuThich.DELETE("", controllers.DeleteYeuThich)
	}
}
