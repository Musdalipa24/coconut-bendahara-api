package model

import "time"

type Pemasukan struct {
	Id          string
	Tanggal     time.Time
	Kategori    string
	Keterangan  string
	Nominal     uint64
	Nota        string
	IdTransaksi string
}
