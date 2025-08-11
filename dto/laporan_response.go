package dto

type LaporanKeuanganResponse struct {
	Id          string `json:"id_laporan"`
	Tanggal     string `json:"tanggal"`
	Keterangan  string `json:"keterangan"`
	Pemasukan   uint64 `json:"pemasukan"`
	Pengeluaran uint64 `json:"pengeluaran"`
	Saldo       int64  `json:"total_saldo"`
}
