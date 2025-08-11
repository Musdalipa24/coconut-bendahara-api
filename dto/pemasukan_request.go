package dto

type PemasukanRequest struct {
	Tanggal    string `json:"tanggal"`
	Kategori   string `json:"kategori"`
	Keterangan string `json:"keterangan"`
	Nominal    uint64 `json:"nominal"`
	Nota       string `json:"nota"`
}
