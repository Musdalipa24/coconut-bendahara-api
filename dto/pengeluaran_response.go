package dto

type PengeluaranResponse struct {
	Id         string `json:"id_pengeluaran"`
	Tanggal    string `json:"tanggal"`
	Nota       string `json:"nota"`
	Nominal    uint64 `json:"nominal"`
	Keterangan string `json:"keterangan"`
}

type PengeluaranPaginationResponse struct {
	Items      []PengeluaranResponse `json:"items"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalItems int                   `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
}
