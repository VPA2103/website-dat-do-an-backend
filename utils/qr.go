package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/skip2/go-qrcode"
	"github.com/subiz/vietqr"
)

func GenerateQRBytes(url string) ([]byte, error) {
	return qrcode.Encode(url, qrcode.Medium, 256)
}

//amount := 12000.0
//bankBIN := "970423"
//account := "00005897596"
//note := "Ủng hộ lũ lụt"

func GenerateQRPayment(amount float64, bankBIN string, accountnumber, note string) (string, error) {
	code := vietqr.Generate(amount, bankBIN, accountnumber, note)

	png, err := qrcode.Encode(code, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Chuyển []byte → base64 string
	base64Img := base64.StdEncoding.EncodeToString(png)

	return base64Img, nil
}

func GenerateSePayQR(account string, bank string, amount int, content string) string {

	// encode nội dung chuyển khoản
	encodedContent := url.QueryEscape(content)

	qrURL := fmt.Sprintf(
		"https://qr.sepay.vn/img?acc=%s&bank=%s&amount=%d&des=%s",
		account,
		bank,
		amount,
		encodedContent,
	)

	return qrURL
}