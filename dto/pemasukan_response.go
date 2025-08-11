package dto

type PemasukanResponse struct {
	Id         string `json:"id_pemasukan"`
	Tanggal    string `json:"tanggal"`
	Kategori   string `json:"kategori"`
	Keterangan string `json:"keterangan"`
	Nominal    uint64 `json:"nominal"`
	Nota       string `json:"nota"`
}

type PemasukanPaginationResponse struct {
	Items      []PemasukanResponse `json:"items"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalItems int                 `json:"total_items"`
	TotalPages int                 `json:"total_pages"`
}
