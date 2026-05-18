package utils

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("MY_SECRET_KEY")

func ParseToken(r *http.Request) (uint, string, error) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		tokenStr = r.Header.Get("Authorization")
	}

	if tokenStr == "" {
		return 0, "", errors.New("missing token")
	}

	// ✅ Đọc từ .env thay vì hardcode
	secretKey := os.Getenv("SECRET_KEY")
	log.Println("🔑 SECRET_KEY:", secretKey)
	if secretKey == "" {
		secretKey = "SECRET_KEY"
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil // ✅ dùng secretKey từ .env
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}

	userID := uint(claims["id"].(float64))
	role := claims["role"].(string)

	return userID, role, nil
}

type JWTClaims struct {
	UserID   uint   `json:"ma_nv"`
	Username string `json:"username"`
	Role     string `json:"role"` // 🔥 thêm dòng này
	jwt.RegisteredClaims
}

func GenerateToken(id uint, email string, role string) (string, error) {
	expireTime := time.Now().Add(3 * time.Hour)
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	}

	fmt.Println("Hết hạn lúc:", expireTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret())
}

// ✅ Hàm dùng trong middleware để xác thực token
func SecretKey() []byte {
	return secretKey
}

func GenerateInvoice() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf(
		"INV-%d-%04d",
		time.Now().Unix(),
		rand.Intn(10000),
	)
}
