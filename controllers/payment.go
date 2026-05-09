package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/utils"
)

type QRRequest struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}

func GetVietQR(c *gin.Context) {
	var req QRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	qr, err := utils.GenerateQRPayment(
		req.Amount,
		config.PaymentCfg.BankBin,
		config.PaymentCfg.AccountNo,
		req.Note,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tạo được QR: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		//"bank_bin":   config.PaymentCfg.BankBin,
		//"account_no": config.PaymentCfg.AccountNo,
		"qr_base64": "data:image/png;base64," + qr,
	})
}
