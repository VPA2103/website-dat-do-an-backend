
package models

import "time"
type ThreadMessage struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	ThreadID string `gorm:"index" json:"thread_id"`

	Role    string `json:"role"`    // user | assistant
	Content string `json:"content"`

	CreatedAt time.Time `json:"created_at"`

	Thread Thread `gorm:"foreignKey:ThreadID;references:ID" json:"-"`
}