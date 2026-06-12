package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func MonAnRoutes(r *gin.Engine,h *controllers.ChatHandler) {
	mon_an := r.Group("/mon_an")
	{
		mon_an.GET("/get-mon-an-co-binh-luan-va-danh-gia", controllers.GetMonAnCoBinhLuanVaDanhGia)
		mon_an.GET("/get-mon-an-co-binh-luan-va-danh-gia-cua-nguoi-dung",middleware.AuthMiddleware(), controllers.GetMonAnCoBinhLuanVaDanhGiaCuaNguoiDung)
		mon_an.POST("/create", h.CreateMonAn)
		mon_an.GET("/all", controllers.GetAllMonAn)
		mon_an.GET("/:id", controllers.GetMonAnByID)
		mon_an.GET("/:id/detail", controllers.GetMonAnDetail)
		mon_an.PATCH("/update/:id", h.UpdateMonAn)
		mon_an.DELETE("/delete/:id", controllers.DeleteMonAn)
		mon_an.GET("/search", controllers.SearchMonAn)
		
	}
}
