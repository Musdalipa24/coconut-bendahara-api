package dto

type TransactionResponse struct {
	Id             string `json:"id_transaksi"`
	Tanggal        string `json:"tanggal"`
	Keterangan     string `json:"keterangan"`
	JenisTransaksi string `json:"jenis_transaksi"`
	Nominal        uint64 `json:"nominal"`
}
