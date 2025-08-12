package dto

type IuranRequest struct {
	Nama        string  `json:"nama" binding:"required"`
	Jumlah      float64 `json:"jumlah" binding:"required"`
	Tanggal     string  `json:"tanggal" binding:"required"` // format: YYYY-MM-DD
	Keterangan  string  `json:"keterangan"`
}


type IuranResponse struct {
    ID         string `json:"id"`
    Nama       string `json:"nama"`
    Jumlah     int    `json:"jumlah"`
    Tanggal    string `json:"tanggal"`
    Keterangan string `json:"keterangan"`
    Message    string `json:"message,omitempty"`
}
