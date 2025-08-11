package model

import "time"

type LaporanKeuangan struct {
	Id          string
	Tanggal     time.Time
	Keterangan  string
	Pemasukan   uint64
	Pengeluaran uint64
	Saldo       int64
}
