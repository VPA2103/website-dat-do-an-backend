package send_mail

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/utils"
)

func SendMailAPI(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := SendDatBanMail(req.Email); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Gửi mail thành công"})
}

func SendDatBanMail(email string) error {
	return utils.SendMail(
		email,
		"Xác nhận đặt bàn",
		"<h3>Đặt bàn thành công</h3>"+
			"<p>Hẹn gặp bạn!</p>",
	)
}


