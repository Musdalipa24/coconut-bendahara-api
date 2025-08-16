package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type LaporanKeuanganService interface {
	GetAllLaporan(ctx context.Context) ([]dto.LaporanKeuanganResponse, int, error)
	GetLastBalance(ctx context.Context) (int64, int, error)
	GetTotalIncome(ctx context.Context) (uint64, int, error)
	GetTotalExpenditure(ctx context.Context) (uint64, int, error)
	GetLaporanByDateRange(ctx context.Context, startDate string, endDate string) ([]dto.LaporanKeuanganResponse, int, error)
}

type laporanKeuanganServiceImpl struct {
	LaporanRepo repository.LaporanKeuanganRepo
	DB          *sql.DB
}

func NewLaporanKeuanganService(laporanRepo repository.LaporanKeuanganRepo, db *sql.DB) LaporanKeuanganService {
	return &laporanKeuanganServiceImpl{
		LaporanRepo: laporanRepo,
		DB:          db,
	}
}

// GetAllLaporan implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetAllLaporan(ctx context.Context) ([]dto.LaporanKeuanganResponse, int, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	laporan, err := l.LaporanRepo.GetAllLaporan(ctx, tx)
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get financial statements")
	}

	return util.ConvertLaporanToListResponseDTO(laporan), http.StatusOK, nil
}

// GetLastBalance implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetLastBalance(ctx context.Context) (int64, int, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	saldo, err := l.LaporanRepo.GetLastBalance(ctx, tx)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("failed to get last balance: %v", err)
	}

	return saldo, http.StatusOK, nil
}

// GetTotalExpenditure implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetTotalExpenditure(ctx context.Context) (uint64, int, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	laporan, err := l.LaporanRepo.GetAllLaporan(ctx, tx)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("failed to get financial statements")
	}

	var totalPengeluaran uint64
	for _, data := range laporan {
		totalPengeluaran += data.Pengeluaran
	}

	return totalPengeluaran, http.StatusOK, nil
}

// GetTotalIncome implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetTotalIncome(ctx context.Context) (uint64, int, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	laporan, err := l.LaporanRepo.GetAllLaporan(ctx, tx)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("failed to get financial statements")
	}

	var totalPemasukan uint64
	for _, data := range laporan {
		totalPemasukan += data.Pemasukan
	}

	return totalPemasukan, http.StatusOK, nil
}

func (l *laporanKeuanganServiceImpl) GetLaporanByDateRange(ctx context.Context, startDate string, endDate string) ([]dto.LaporanKeuanganResponse, int, error) {
	// Mulai transaksi
	tx, err := l.DB.Begin()
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	// Panggil repository untuk mendapatkan data laporan berdasarkan rentang tanggal
	laporans, err := l.LaporanRepo.GetLaporanByDateRange(ctx, tx, startDate, endDate)
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get laporan by date range: %v", err)
	}

	// Konversi data laporan ke DTO
	return util.ConvertLaporanToListResponseDTO(laporans), http.StatusOK, nil
}
