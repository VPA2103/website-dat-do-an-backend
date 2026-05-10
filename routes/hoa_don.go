package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func HoaDonRoutes(r *gin.Engine) {

	hoadon := r.Group("/hoa-don")
	{
		hoadon.POST("/dat-do-an", controllers.DatDoAn)
	}
}
