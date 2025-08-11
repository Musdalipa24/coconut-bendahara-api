package model

import "time"

type Iuran struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama_warga"`
	Jumlah    float64   `json:"jumlah"`
	Status    string    `json:"status"` // Lunas, Belum Lunas
	CreatedAt time.Time `json:"created_at"`
}
