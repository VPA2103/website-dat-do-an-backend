package models

import (
	"time"

	"gorm.io/datatypes"
)

type Payment struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement;type:bigint"`
	OrderID       uint64         `gorm:"type:bigint;index"` // Đánh index để query theo đơn hàng nhanh hơn
	Provider      string         `gorm:"type:varchar(50)"`  // VD: sepay
	InvoiceNumber string         `gorm:"type:varchar(100);unique"`
	Amount        float64        `gorm:"type:decimal(15,2)"`
	Status        string         `gorm:"type:varchar(20);default:'pending'"` // pending, paid, failed
	TransactionID string         `gorm:"type:varchar(255)"` // Mã giao dịch từ SePay
	BankCode      string         `gorm:"type:varchar(50)"`
	RawData       datatypes.JSON `gorm:"type:json"` // Lưu payload webhook để debug
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
