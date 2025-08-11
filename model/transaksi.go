package model

import "time"

type Transaction struct {
	Id             string
	Tanggal        time.Time
	Keterangan     string
	JenisTransaksi string
	Nominal        uint64
}
