package usecase

import "github.com/vpa/quanlynhahang-backend/internal/dto"

type RealtimeSender interface {
	Broadcast( msg dto.WSMessage)
	SendToUser(userID uint, msg dto.WSMessage)
	SendToRole(role string, msg dto.WSMessage)
}
