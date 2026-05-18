package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
    fmt.Println("🔍 Data to sign:", data)   // Debug

    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(data))
    sig := base64.StdEncoding.EncodeToString(h.Sum(nil))

    fmt.Println("🔑 Generated signature:", sig)  // Debug
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

func SePayWebhookHandler(c *gin.Context) {
	// Đọc raw body (quan trọng nếu sau này làm HMAC)
	// bodyBytes, err := c.GetRawData()

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
	// 	return
	// }

	// Parse JSON
	var payload SepayWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil { // hoặc json.Unmarshal(bodyBytes, &payload)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// === TODO: Xử lý business logic ===
	// 1. Kiểm tra HMAC (nếu bạn bật xác thực)
	// 2. Kiểm tra idempotency với payload.ID
	// 3. Cập nhật trạng thái đơn hàng, ghi log, etc.

	log.Printf("[SePay Webhook] Received - ID: %d | Amount: %d VND | Content: %s | Code: %v",
		payload.ID, payload.TransferAmount, payload.Content, payload.Code)

	// Phản hồi BẮT BUỘC theo yêu cầu của SePay
	c.JSON(http.StatusOK, gin.H{"success": true})
}
