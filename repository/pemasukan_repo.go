package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type PemasukanRepo interface {
	AddPemasukan(ctx context.Context, tx *sql.Tx, pemasukan model.Pemasukan) (model.Pemasukan, error)
	UpdatePemasukan(ctx context.Context, tx *sql.Tx, pemasukan model.Pemasukan, no string) (model.Pemasukan, error)
	GetPemasukan(ctx context.Context, tx *sql.Tx, page int, pageSize int) ([]model.Pemasukan, int, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (model.Pemasukan, error)
	DeletePemasukan(ctx context.Context, tx *sql.Tx, pemasukan model.Pemasukan) (model.Pemasukan, error)
	GetPemasukanByDateRange(ctx context.Context, tx *sql.Tx, startDate, endDate string, page int, pageSize int) ([]model.Pemasukan, int, error)
}