package models

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoomID    uint      `gorm:"index:idx_room_created" json:"room_id"`
	SenderID  uint      `gorm:"index" json:"sender_id"`
	Content   string    `gorm:"type:text" json:"content"`
	Type      string    `gorm:"size:20;default:'text'" json:"type"`
	CreatedAt time.Time `gorm:"index:idx_room_created" json:"created_at"`

	Room   Room     `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Sender NguoiDung `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

type MessageReadReceipt struct {
	MessageID uint      `gorm:"primaryKey" json:"message_id"`
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	ReadAt    time.Time `gorm:"index" json:"read_at"`
}
