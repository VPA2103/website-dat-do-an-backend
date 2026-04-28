package models

import (
	"time"

	"gorm.io/datatypes"
)

type Payments struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	OrderID       uint64         `gorm:"index"`
	Provider      string         `gorm:"type:varchar(50)"`
	InvoiceNumber string         `gorm:"type:varchar(100);unique"`
	Amount        float64        `gorm:"type:decimal(15,2)"`
	Status        string         `gorm:"type:varchar(20);default:'pending'"`
	TransactionID string         `gorm:"type:varchar(255)"`
	BankCode      string         `gorm:"type:varchar(50)"`
	RawData       datatypes.JSON `gorm:"type:jsonb"` // PostgreSQL nên dùng jsonb

	CreatedAt time.Time
	UpdatedAt time.Time
}
