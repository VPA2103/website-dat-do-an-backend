package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func DatBanRoutes(r *gin.Engine) {
	datban := r.Group("/dat-ban")
	{
		// Khách
		datban.POST("", middleware.AuthMiddleware(), controllers.CreateDatBan)                                    // tạo đặt bàn //ok
		datban.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetAllDatBan) // danh sách //ok
		datban.GET("/:id", middleware.AuthMiddleware(), controllers.GetDatBanByID)                                // chi tiết
		datban.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateDatBan)                                 // sửa thông tin
		datban.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteDatBan)                              //ok

		// Nhân viên
		datban.PUT("/:id/xac-nhan", middleware.AuthMiddleware(), controllers.XacNhanDatBan)
		datban.GET("/lay-danh-sach-dat-ban-cua-nguoi-dung", middleware.AuthMiddleware(), controllers.GetDatBanCuaNguoiDung)
		datban.PUT("/huy-dat-ban/:id", middleware.AuthMiddleware(), controllers.HuyDatBan)
	}
}
