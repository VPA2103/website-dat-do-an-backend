package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
	"gorm.io/datatypes"
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

	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

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
		Provider:      "sepay",
		InvoiceNumber: invoice,
		Amount:        order.TongTien,
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
	baseURL := "https://pgapi-sandbox.sepay.vn/payment/link"
	//Production	https://pgapi.sepay.vn
	//Sandbox	https://pgapi-sandbox.sepay.vn

	merchant := os.Getenv("SEPAY_MERCHANT_ID")
	secret := os.Getenv("SEPAY_SECRET_KEY")

	params := map[string]string{
		"merchant":             merchant,
		"order_amount":         fmt.Sprintf("%d", amount),
		"currency":             "VND",
		"operation":            "PURCHASE",
		"order_description":    "Thanh toán đơn hàng",
		"order_invoice_number": invoice,
		"success_url":          "https://desirous-rodger-panlogistically.ngrok-free.dev/success",
		"cancel_url":           "https://desirous-rodger-panlogistically.ngrok-free.dev/cancel",
		"error_url":            "https://desirous-rodger-panlogistically.ngrok-free.dev/error",
	}

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

	query := buildQuery(params, keys)

	fmt.Println("STRING TO SIGN:", query)

	signature := sign(query, secret)

	fmt.Println("SIGNATURE:", signature)

	// build final URL
	v := url.Values{}

	for _, k := range keys {
		v.Set(k, params[k])
	}

	// QUAN TRỌNG
	v.Set("signature", signature)

	finalURL := baseURL + "?" + v.Encode()

	fmt.Println("FINAL URL:", finalURL)

	return finalURL
}

func HandleIPN(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "cannot read body"})
		return
	}
	fmt.Println("SEPAY WEBHOOK:", string(body))
	signature := c.GetHeader("X-SePay-Signature")

	// verify signature trước
	if !verifySignature(body, signature) {
		c.JSON(400, gin.H{"error": "invalid signature"})
		return
	}

	// parse JSON 1 lần duy nhất
	var payload SePayWebhook
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	invoice := payload.Order.InvoiceNumber

	// 1. find payment
	var payment models.Payments
	if payment.Status == "paid" {
		return
	}
	if err := config.DB.
		Where("invoice_number = ?", invoice).
		First(&payment).Error; err != nil {
		c.Status(200)
		return
	}

	// 2. idempotent
	if payment.Status == "paid" {
		c.Status(200)
		return
	}

	// 3. verify amount (SAU KHI LẤY PAYMENT)
	if float64(payload.Order.Amount) < payment.Amount {
		c.JSON(400, gin.H{"error": "invalid amount"})
		return
	}

	// 4. update payment + raw data
	if payload.NotificationType == "ORDER_PAID" {

		payment.Status = "paid"
		payment.RawData = datatypes.JSON(body)
		payment.TransactionID = payload.Order.InvoiceNumber // nếu API có txn id thì thay

		config.DB.Save(&payment)

		// 5. update order (KHÔNG đổi DB structure)
		config.DB.Model(&models.HoaDon{}).
			Where("ma_hd = ?", payment.OrderID).
			Update("trang_thai", "da_thanh_toan")
	}

	c.Status(200)
}
