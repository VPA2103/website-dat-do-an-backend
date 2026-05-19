package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func MonAnRoutes(r *gin.Engine) {
	mon_an := r.Group("/mon_an")
	{
		mon_an.POST("/create", controllers.CreateMonAn)
		mon_an.GET("/all", controllers.GetAllMonAn)
		mon_an.GET("/:id", controllers.GetMonAnByID)
		mon_an.GET("/:id/detail", controllers.GetMonAnDetail)
		mon_an.PATCH("/update/:id", controllers.UpdateMonAn)
		mon_an.DELETE("/delete/:id", controllers.DeleteMonAn)
	}
}
