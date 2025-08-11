package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type PengeluaranRepo interface {
	AddPengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran) (model.Pengeluaran, error)
	UpdatePengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran, no string) (model.Pengeluaran, error)
	GetPengeluaran(ctx context.Context, tx *sql.Tx, page int, pageSize int) ([]model.Pengeluaran, int, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (model.Pengeluaran, error)
	DeletePengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran) (model.Pengeluaran, error)
	GetPengeluaranByDateRange(ctx context.Context, tx *sql.Tx, startDate, endDate string, page int, pageSize int) ([]model.Pengeluaran, int, error)
}
