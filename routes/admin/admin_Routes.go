package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func AdminRoutes(r *gin.Engine) {
	nhanvien := r.Group("/admin")
	{
		// ✅ Chỉ admin được phép
		nhanvien.POST("/create/nhanvien", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.CreateNhanVien)
		nhanvien.PATCH("/update/nhanvien/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.UpdateNhanVien)
		nhanvien.DELETE("/delete/nhanvien/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.DeleteNhanVien)
		// ✅ chỉ admin có thể xem danh sách
		nhanvien.GET("/nhanvien/layTatCa", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetAllNhanVien)

	}
}
