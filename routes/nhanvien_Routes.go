package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func NguoiDungRoutes(r *gin.Engine, hub *websocket.Hub) {
	nguoidung := r.Group("/nhanvien")
	{
		// ✅ Chỉ admin được phép
		nguoidung.POST("/create", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.CreateNhanVien)
		nguoidung.PATCH("/update/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.UpdateNhanVien)
		nguoidung.DELETE("/delete/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.DeleteNhanVien)
		nguoidung.POST("/create-nv-db", controllers.CreateNhanVien)

		nguoidung.GET("/layRaThongTinNhanVien/:id", controllers.GetNhanVienByID)

		// ✅ Chỉ nhân viên được phép
		nguoidung.PATCH("/capNhatThongTinCaNhan/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"), controllers.UpdateThongTinCaNhan)

		// ✅ Cả admin và user đều có thể xem danh sách
		nguoidung.GET("/layTatCa", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "user"), controllers.GetAllNhanVien)

		nguoidung.POST("/assign-shipper", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.AssignShipper(hub))

	}
}
