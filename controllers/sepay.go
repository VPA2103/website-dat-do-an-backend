package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
)

// Sandbox: "https://pgapi-sandbox.sepay.vn"
//https://pgapi.sepay.vn

const (
	SePayMerchantID = "SP-TEST-VP634685"
	SePaySecretKey  = "spsk_test_gxUG9RkZiMW2DtP7EweSYHqQeXtMyXF5"
	// Endpoint JSON
	SePayCheckoutAPI = "https://pay.sepay.vn/v1/checkout/init"
)

type CheckoutRequest struct {
	Amount        int64  `json:"amount" binding:"required"`
	InvoiceNumber string `json:"invoice_number" binding:"required"`
	Description   string `json:"description" binding:"required"`
	SuccessURL    string `json:"success_url"`
	ErrorURL      string `json:"error_url"`
	CancelURL     string `json:"cancel_url"`
}

type SePayResponse struct {
	Message     string `json:"message"`
	RedirectURL string `json:"redirect_url"`
}

func generateSignature(fields map[string]string, secretKey string) string {
	// Thứ tự fields quan trọng - theo docs SePay
	order := []string{
		"order_amount",
		"merchant",
		"currency",
		"operation",
		"order_description",
		"order_invoice_number",
		"success_url",
		"error_url",
		"cancel_url",
	}

	var parts []string
	for _, key := range order {
		if val, exists := fields[key]; exists && val != "" {
			parts = append(parts, key+"="+val)
		}
	}

	data := strings.Join(parts, ",")
	fmt.Println("🔍 Data to sign:", data) // Debug

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	sig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	fmt.Println("🔑 Generated signature:", sig) // Debug
	return sig
}

func CreateSePayPaymentForm(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fields := map[string]string{
		"merchant":             SePayMerchantID,
		"currency":             "VND",
		"order_amount":         fmt.Sprintf("%d", req.Amount),
		"operation":            "PURCHASE",
		"order_description":    req.Description,
		"order_invoice_number": req.InvoiceNumber,
		"success_url":          req.SuccessURL,
		"error_url":            req.ErrorURL,
		"cancel_url":           req.CancelURL,
	}

	fields["signature"] = generateSignature(fields, SePaySecretKey)

	// Tạo form HTML
	html := `<form action="https://pay.sepay.vn/v1/checkout/init" method="POST" id="sepayForm">`
	for k, v := range fields {
		html += fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, k, template.HTMLEscapeString(v))
	}
	html += `</form>`
	html += `<script>document.getElementById("sepayForm").submit();</script>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

type SepayWebhookPayload struct {
	ID              int64   `json:"id"`
	Gateway         string  `json:"gateway"`
	TransactionDate string  `json:"transactionDate"`
	AccountNumber   string  `json:"accountNumber"`
	SubAccount      string  `json:"subAccount"`
	Code            *string `json:"code"` // có thể null
	Content         string  `json:"content"`
	TransferType    string  `json:"transferType"`
	Description     string  `json:"description"`
	TransferAmount  int64   `json:"transferAmount"`
	Accumulated     int64   `json:"accumulated"`
	ReferenceCode   string  `json:"referenceCode"`
}

func SePayWebhook(c *gin.Context) {

	var payload struct {
		Content string  `json:"content"`
		Amount  float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	// content: HD25
	if payload.Content == "SEPAY TEST WEBHOOK" {
		c.JSON(200, gin.H{
			"message": "Webhook test received",
		})
		return
	}

	var hoaDon models.HoaDon

	err := config.DB.
		Where("ma_hd = ?", strings.ReplaceAll(payload.Content, "HD", "")).
		First(&hoaDon).Error

	if err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy hóa đơn"})
		return
	}

	// check đủ tiền
	if payload.Amount < hoaDon.TongTien {

		c.JSON(400, gin.H{
			"error": "Thanh toán thiếu tiền",
		})
		return
	}

	config.DB.Model(&hoaDon).Updates(map[string]interface{}{
		"trang_thai_thanh_toan": "da_thanh_toan",
	})

	c.JSON(200, gin.H{
		"success": true,
	})
}

func GetQR(c *gin.Context) {

	qr := utils.GenerateSePayQR(
		"0123456789",
		"MBBank",
		50000,
		"DON123",
	)

	c.JSON(200, gin.H{
		"qr_url": qr,
	})
}
