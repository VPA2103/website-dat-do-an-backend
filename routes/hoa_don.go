package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func HoaDonRoutes(r *gin.Engine, hub *websocket.Hub) {

	ctrl := controllers.NewHoaDonController(hub)

	hoaDon := r.Group("/hoa-don")
	{
		hoaDon.POST("", middleware.AuthMiddleware(), ctrl.DatDoAn)

		hoaDon.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "shipper"), ctrl.GetHoaDons)

		hoaDon.PUT("/:id", ctrl.UpdateHoaDon)

		hoaDon.PUT("/:id/trang-thai", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "shipper"), ctrl.UpdateTrangThaiHoaDon)

		hoaDon.GET("/user", middleware.AuthMiddleware(), ctrl.GetHoaDonByNguoiDung)

		hoaDon.GET("/trang-thai", ctrl.GetHoaDonByTrangThai)

		hoaDon.GET("/:id", middleware.AuthMiddleware(), ctrl.GetHoaDonByID)

		hoaDon.PUT("/:id/huy", ctrl.HuyHoaDon)

		hoaDon.DELETE("/:id", ctrl.XoaHoaDon)

		hoaDon.PUT(":id/huy_thanh_toan", middleware.AuthMiddleware(), ctrl.HuyHoaDonNguoiDung)

		hoaDon.GET("/cho-thanh-toan", middleware.AuthMiddleware(), ctrl.GetHoaDonChoThanhToan)
		//thongke
		hoaDon.GET("/doanh-thu-ngay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.GetDoanhThuTheoNgay)
		hoaDon.GET("/doanh-thu-thang", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.GetDoanhThuTheoThang)
		hoaDon.GET("/doanh-thu-nam", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.GetDoanhThuTheoNam)
		hoaDon.GET("/mon-an-ban-chay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.GetTopMonAnBanChay)
		hoaDon.GET("/ti-le-hoan-thanh-hom-nay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.GetTiLeHoanThanhHomNay)
		hoaDon.GET("/so-don-theo-ngay", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.SoDonTheoNgay)
		hoaDon.GET("/don-da-giao-hom-nay",middleware.AuthMiddleware(),middleware.RoleMiddleware("admin"),ctrl.DonHangDaGiaoHomNay)
		hoaDon.GET("/top-mon-ban-chay-nhat", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), ctrl.TopMonBanChay)
		// hoaDon.POST("/:ma_hd/thanh-toan", controllers.ThanhToanHoaDon)
	}
}
