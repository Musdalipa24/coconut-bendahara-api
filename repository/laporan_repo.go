package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type LaporanKeuanganRepo interface {
	GetAllLaporan(ctx context.Context, tx *sql.Tx) ([]model.LaporanKeuangan, error)
	GetLastBalance(ctx context.Context, tx *sql.Tx) (int64, error)
	GetLaporanByDateRange(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]model.LaporanKeuangan, error)
}