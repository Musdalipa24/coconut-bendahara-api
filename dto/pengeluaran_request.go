package dto

type PengeluaranRequest struct {
	Tanggal    string `json:"tanggal"`
	Nota       string `json:"nota"`
	Nominal    uint64 `json:"nominal"`
	Keterangan string `json:"keterangan"`
}
