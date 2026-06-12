package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func AdminRoutes(r *gin.Engine) {

	admin := r.Group("/admin")
	{
		// ✅ Chỉ admin được phép
		admin.POST("/create/nhanvien", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.CreateNhanVien)
		admin.GET("/doanh-thu-ngay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetDoanhThuTheoNgay)
		admin.GET("/doanh-thu-thang", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetDoanhThuTheoThang)
		admin.GET("/doanh-thu-nam", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetDoanhThuTheoNam)
		admin.GET("/export-doanh-thu-ngay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.ExportDoanhThuNgay)
		admin.GET("/danh-sach-doanh-thu-ngay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetDanhSachNgayDoanhThu)
		admin.GET("/top-mon-ban-chay-nhat", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.TopMonBanChay)
		admin.GET("/ti-le-hoan-thanh-hom-nay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetTiLeHoanThanhHomNay)
		admin.GET("/mon-an-ban-chay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetTopMonAnBanChay)
	}
}
