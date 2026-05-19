package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func OptionRoutes(r *gin.Engine) {
	option := r.Group("/option")
	{
		// nhóm option
		option.POST("/nhom", controllers.CreateNhomOption)
		option.GET("/nhom", controllers.GetAllNhomOption)
		option.GET("/nhom/:id", controllers.GetNhomOptionByID)
		option.PUT("/nhom/:id", controllers.UpdateNhomOption)
		option.DELETE("/nhom/:id", controllers.DeleteNhomOption)

		// option item
		option.POST("/item", controllers.CreateOptionItem)
		option.GET("/item", controllers.GetAllOptionItem)
		option.GET("/item/:id", controllers.GetOptionItemByID)
		option.PUT("/item/:id", controllers.UpdateOptionItem)
		option.DELETE("/item/:id", controllers.DeleteOptionItem)
	}
}
