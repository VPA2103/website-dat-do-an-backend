package dto

type DiaChiInput struct {
	HoTen    string  `json:"ho_ten" binding:"required"`
	SDT      string  `json:"sdt" binding:"required"`
	DiaChi   string  `json:"dia_chi" binding:"required"`
	Latitude float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	MacDinh  bool    `json:"mac_dinh"`
}