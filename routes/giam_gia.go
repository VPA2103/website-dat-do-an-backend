package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func GiamGiaRoutes(r *gin.Engine) {
	giamGia := r.Group("/giam-gia")
	{
		giamGia.POST("/", controllers.CreateGiamGia)
		giamGia.GET("/", controllers.GetAllGiamGia)
		giamGia.GET("/:id", controllers.GetGiamGiaById)
		giamGia.PATCH("/:id", controllers.UpdateGiamGia)
		giamGia.DELETE("/:id", controllers.DeleteGiamGia)
	}
}
