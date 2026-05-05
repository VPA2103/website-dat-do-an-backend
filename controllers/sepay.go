package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

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

func verifySignature(body []byte, signature string) bool {
	secret := os.Getenv("SEPAY_SECRET_KEY")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func CreatePaymentService(orderID int64) (map[string]interface{}, error) {
	invoice := utils.GenerateInvoice()

	var order models.HoaDon
	if err := config.DB.First(&order, orderID).Error; err != nil {
		return nil, err
	}

	payment := models.Payments{
		OrderID:       uint64(orderID),
		InvoiceNumber: invoice,
		Status:        "pending",
	}
	config.DB.Create(&payment)

	url := buildSePayURL(invoice, int64(order.TongTien))

	return map[string]interface{}{
		"checkout_url": url,
	}, nil
}

func CreatePayment(c *gin.Context) {
	var req struct {
		OrderID int64 `json:"order_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	data, err := CreatePaymentService(req.OrderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}

func sign(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func buildQuery(params map[string]string, keys []string) string {
	var buf strings.Builder

	for i, k := range keys {
		buf.WriteString(k + "=" + params[k])
		if i < len(keys)-1 {
			buf.WriteString("&")
		}
	}

	return buf.String()
}

func buildSePayURL(invoice string, amount int64) string {
	baseURL := "https://pay-sandbox.sepay.vn/v1/checkout"
	secret := os.Getenv("SEPAY_SECRET_KEY")

	params := url.Values{}
	params.Set("merchant", os.Getenv("SEPAY_MERCHANT_ID"))
	params.Set("order_amount", fmt.Sprintf("%d", amount))
	params.Set("currency", "VND")
	params.Set("operation", "PURCHASE")
	params.Set("order_description", "Thanh toan don hang")
	params.Set("order_invoice_number", invoice)
	params.Set("success_url", "https://desirous-rodger-panlogistically.ngrok-free.dev/success")
	params.Set("cancel_url", "https://desirous-rodger-panlogistically.ngrok-free.dev/cancel")
	params.Set("error_url", "https://desirous-rodger-panlogistically.ngrok-free.dev/error")

	keys := []string{
		"cancel_url",
		"currency",
		"error_url",
		"merchant",
		"operation",
		"order_amount",
		"order_description",
		"order_invoice_number",
		"success_url",
	}
	// ✅ chuẩn hóa string
	query := buildQuery(map[string]string{
		"cancel_url":           params.Get("cancel_url"),
		"currency":             "VND",
		"error_url":            params.Get("error_url"),
		"merchant":             params.Get("merchant"),
		"operation":            "PURCHASE",
		"order_amount":         params.Get("order_amount"),
		"order_description":    params.Get("order_description"),
		"order_invoice_number": params.Get("order_invoice_number"),
		"success_url":          params.Get("success_url"),
	}, keys)

	fmt.Println("STRING TO SIGN:", query)
	fmt.Println("TO SIGN:", query)
	signature := sign(query, secret)
	fmt.Println("SIGN:", signature)
	fmt.Println("MERCHANT:", os.Getenv("SEPAY_MERCHANT_ID"))
	fmt.Println("SECRET:", secret)
	params.Set("signature", signature)

	return baseURL + "?" + params.Encode()
}

func HandleIPN(c *gin.Context) {

	body, _ := io.ReadAll(c.Request.Body)

	signature := c.GetHeader("X-SePay-Signature")

	if !verifySignature(body, signature) {
		c.JSON(400, gin.H{"error": "invalid signature"})
		return
	}

	var payload SePayWebhook

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	invoice := payload.Order.InvoiceNumber

	var payment models.Payments
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
