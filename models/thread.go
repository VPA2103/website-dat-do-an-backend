package models

import "time"

type Thread struct {
	ID string `gorm:"primaryKey" json:"id"`

	CreatedAt time.Time `json:"created_at"`

	Messages []ThreadMessage `gorm:"foreignKey:ThreadID" json:"messages"`
}