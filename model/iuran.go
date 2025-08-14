package model

import "database/sql"

type Iuran struct {
	IdIuran  sql.NullString
	Periode  sql.NullString
	MingguKe sql.NullInt64
}

type PembayaranIuran struct {
	IdPembayaran sql.NullString
	IdMember     sql.NullString
	Status       sql.NullString
	TanggalBayar sql.NullString
	Iuran        Iuran
}

type Member struct {
	IdMember         string
	NRA              string
	Nama             string
	Status           string
	CreatedAt        string
	UpdatedAt        string
	PembayaranIurans []PembayaranIuran
}
