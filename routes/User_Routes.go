package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func UserRoutes(r *gin.Engine) {

	User := r.Group("/user")
	{
		User.GET("/so-tien-da-mua", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"), controllers.GetTongTienDaMua)


	}
}
