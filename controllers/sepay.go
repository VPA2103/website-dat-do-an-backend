package controllers

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
)

type SePayWebhook struct {
	NotificationType string `json:"notification_type"`
	Order            struct {
		InvoiceNumber string `json:"invoice_number"`
		Amount        int64  `json:"amount"`
		Status        string `json:"status"`
	} `json:"order"`
}

func CreatePayment(orderID int64) (string, error) {
	invoice := utils.GenerateInvoice()

	// 👉 lấy order từ DB để có amount
	var order models.HoaDon

	if err := config.DB.First(&order, orderID).Error; err != nil {
		return "", err
	}

	payment := models.Payment{
		OrderID:       uint64(orderID),
		InvoiceNumber: invoice,
		Status:        "pending",
	}

	config.DB.Create(&payment)

	// ✅ truyền đủ 2 tham số
	paymentURL := buildSePayURL(invoice, order.TongTien)

	return paymentURL, nil
}

func buildSePayURL(invoice string, amount float64) string {
	base := "https://sandbox.sepay.vn/payment"

	params := url.Values{}
	params.Add("merchant_id", "YOUR_MERCHANT_ID")
	params.Add("invoice", invoice)
	params.Add("amount", fmt.Sprintf("%.0f", amount))
	params.Add("return_url", "http://localhost:3000/payment/success")

	return base + "?" + params.Encode()
}

func HandleIPN(c *gin.Context) {

	var payload SePayWebhook

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	invoice := payload.Order.InvoiceNumber

	var payment models.Payment
	if err := config.DB.Where("invoice_number = ?", invoice).First(&payment).Error; err != nil {
		c.Status(200)
		return
	}

	// idempotent (tránh update 2 lần)
	if payment.Status == "paid" {
		c.Status(200)
		return
	}

	if payload.NotificationType == "ORDER_PAID" {
		payment.Status = "paid"
		config.DB.Save(&payment)

		// update order
		config.DB.Model(&models.HoaDon{}).
			Where("id = ?", payment.OrderID).
			Update("status", "paid")
	}

	c.Status(200)
}
