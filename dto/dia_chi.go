package dto

type DiaChiInput struct {
	HoTen    string  `json:"ho_ten" binding:"required"`
	SDT      string  `json:"sdt" binding:"required"`
	DiaChi   string  `json:"dia_chi" binding:"required"`
	Latitude float64 `json:"latitude" `
	Longitude float64 `json:"longitude" `
	MacDinh  bool    `json:"mac_dinh"`
}