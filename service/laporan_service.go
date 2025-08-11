package service

import (
	"context"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

type LaporanKeuanganService interface {
	GetAllLaporan(ctx context.Context) ([]dto.LaporanKeuanganResponse, error)
	GetLastBalance(ctx context.Context) (int64, error)
	GetTotalIncome(ctx context.Context) (uint64, error)
	GetTotalExpenditure(ctx context.Context) (uint64, error)
	GetLaporanByDateRange(ctx context.Context, startDate string, endDate string) ([]dto.LaporanKeuanganResponse, error)
}
