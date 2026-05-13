package models

import (
	"time"

	"gorm.io/datatypes"
)

type Payments struct {
	ID            uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID       uint64         `json:"order_id" gorm:"index"`
	Provider      string         `json:"provider" gorm:"type:varchar(50)"`
	InvoiceNumber string         `json:"invoice_number" gorm:"type:varchar(100);unique"`
	Amount        float64        `json:"amount" gorm:"type:decimal(15,2)"`
	Status        string         `json:"status" gorm:"type:varchar(20);default:'pending'"`
	TransactionID string         `json:"transaction_id" gorm:"type:varchar(255)"`
	BankCode      string         `json:"bank_code" gorm:"type:varchar(50)"`
	RawData       datatypes.JSON `json:"raw_data" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
