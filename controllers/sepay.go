package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

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

	return map[string]interface{}{
		"bank_name":      "TPBank",
		"account_number": "00005897596", // STK bạn cấu hình trong SePay
		"account_name":   "Vo Phuc An",
		"amount":         order.TongTien,
		"content":        invoice,
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

// func createSePayPayment(invoice string, amount int64) (string, error) {
// 	url := "https://sepay.vn/api/payment" // ⚠️ có thể sai endpoint

// 	payload := map[string]interface{}{
// 		"merchant_id":    os.Getenv("SEPAY_MERCHANT_ID"),
// 		"amount":         amount,
// 		"invoice_number": invoice,
// 		"return_url":     "https://desirous-rodger-panlogistically.ngrok-free.dev/payment/success",
// 	}

// 	body, _ := json.Marshal(payload)

// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	// ⚠️ có thể sai header (rất nghi ngờ)
// 	req.Header.Set("Authorization", "Bearer "+os.Getenv("SEPAY_SECRET_KEY"))

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	// 🔥 QUAN TRỌNG: đọc raw response
// 	bodyBytes, _ := io.ReadAll(resp.Body)

// 	fmt.Println("STATUS:", resp.Status)
// 	fmt.Println("BODY:", string(bodyBytes))

// 	// thử parse
// 	var result struct {
// 		PaymentURL string `json:"payment_url"`
// 		Error      string `json:"error"`
// 		Message    string `json:"message"`
// 	}

// 	json.Unmarshal(bodyBytes, &result)

// 	// 🔥 nếu không có URL → báo lỗi rõ ràng
// 	if result.PaymentURL == "" {
// 		return "", fmt.Errorf("sepay error: %s | %s | raw: %s",
// 			result.Error, result.Message, string(bodyBytes))
// 	}

// 	return result.PaymentURL, nil
// }

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
