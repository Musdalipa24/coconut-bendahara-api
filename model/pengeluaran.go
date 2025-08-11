package model

import "time"

type Pengeluaran struct {
	Id          string
	Tanggal     time.Time
	Nota        string
	Nominal     uint64
	Keterangan  string
	IdTransaksi string
}
