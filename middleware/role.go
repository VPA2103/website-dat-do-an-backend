package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleAny, exists := c.Get("role")
		if !exists || roleAny == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Không xác định quyền người dùng",
			})
			c.Abort()
			return
		}

		role, ok := roleAny.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Role không hợp lệ",
			})
			c.Abort()
			return
		}

		// kiểm tra role có nằm trong danh sách cho phép không
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Không có quyền truy cập",
		})
		c.Abort()
	}
}
