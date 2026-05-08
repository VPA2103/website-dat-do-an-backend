package usecase

import (
	"github.com/vpa/quanlynhahang-backend/internal/dto"
	"github.com/vpa/quanlynhahang-backend/internal/repository"
	"github.com/vpa/quanlynhahang-backend/models"
)

type NotificationUseCase struct {
	RT   RealtimeSender
	Repo repository.NotificationRepository // thêm
}

func (uc *NotificationUseCase) Notify(msg dto.WSMessage) {
	noti := &models.ThongBao{
		MaNguoiDung:  msg.TargetUserID,
		Title:   "Thông báo",
		Content: msg.Content,
		Type:    msg.Type,
		TrangThai:  false,
	}

	uc.Repo.Create(noti)

	// 🔥 nếu gửi cho admin
	if msg.Role == "admin" {
		uc.RT.SendToRole("admin", msg)
		return
	}

	// gửi user thường
	uc.RT.SendToUser(msg.TargetUserID, msg)
}
