package dto

type WSMessage struct {
	TargetUserID uint   `json:"target_user_id,omitempty"`
	Type         string `json:"type"`
	RoomID       uint   `json:"room_id,omitempty"`
	Content      string `json:"content,omitempty"`
	Role         string `json:"role,omitempty"`
	Payload      interface{} `json:"payload,omitempty"`
}
