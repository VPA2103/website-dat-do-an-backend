package repository

import (
	"github.com/vpa/quanlynhahang-backend/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Save(roomID uint, senderID uint, content string) error
}

// internal/repository/message_repo.go

type GormMessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *GormMessageRepository {
	return &GormMessageRepository{DB: db}
}

func (r *GormMessageRepository) Save(roomID uint, senderID uint, content string) error {
	return r.DB.Create(&models.Message{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}).Error
}

// repo thông báo
type GormNotificationRepository struct {
	DB *gorm.DB
}

type NotificationRepository interface {
	Create(noti *models.ThongBao) error
}

func NewNotificationRepository(db *gorm.DB) *GormNotificationRepository {
	return &GormNotificationRepository{DB: db}
}

func (r *GormNotificationRepository) Create(noti *models.ThongBao) error {
	return r.DB.Create(noti).Error
}
