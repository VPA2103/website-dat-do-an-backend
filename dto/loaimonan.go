package dto

type LoaiMonAn struct {
	MaLoaiMonAn  uint   `json:"ma_loai_mon_an"`
	TenLoaiMonAn string `json:"ten_loai_mon_an" form:"ten_loai_mon_an"`
	AnhLoaiMonAn string `json:"anh_loai_mon_an" form:"anh_loai_mon_an"`
}