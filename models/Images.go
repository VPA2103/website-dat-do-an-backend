package models

type HinhAnh struct {
	ID  uint   `gorm:"primaryKey" json:"id"`
	Url string `json:"url"`
	// Dùng 2 trường để phân biệt loại và id của đối tượng (nhân viên, bàn ăn, v.v.)
	OwnerID   uint   `json:"owner_id"`
	OwnerType string `json:"owner_type"`
}
