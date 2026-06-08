package usecase

import (
	"github.com/vpa/quanlynhahang-backend/dto"
	"github.com/vpa/quanlynhahang-backend/internal/repository"
)

type ChatUseCase struct {
	RT   RealtimeSender
	Repo repository.MessageRepository
}

func (uc *ChatUseCase) SendMessage(userID uint, msg dto.WSMessage) error {
	// 1. Save DB
	if err := uc.Repo.Save(msg.RoomID, userID, msg.Content); err != nil {
		return err
	}

	// 2. Realtime
	uc.RT.Broadcast(msg)

	return nil
}
