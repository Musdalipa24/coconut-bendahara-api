package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type LaporanKeuanganRepo interface {
	GetAllLaporan(ctx context.Context, tx *sql.Tx) ([]model.LaporanKeuangan, error)
	GetLastBalance(ctx context.Context, tx *sql.Tx) (int64, error)
	GetTotalIncome(ctx context.Context, tx *sql.Tx) (uint64, error)
	GetTotalExpenditure(ctx context.Context, tx *sql.Tx) (uint64, error)
	GetLaporanByDateRange(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]model.LaporanKeuangan, error)
}

type laporanKeuanganRepoImpl struct {
}

func NewLaporanKeuanganRepo() LaporanKeuanganRepo {
	return &laporanKeuanganRepoImpl{}
}

// GetAllLaporan implements LaporanKeuanganRepo.
func (l *laporanKeuanganRepoImpl) GetAllLaporan(ctx context.Context, tx *sql.Tx) ([]model.LaporanKeuangan, error) {
	var laporans []model.LaporanKeuangan
	query := "SELECT id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo FROM laporan_keuangan ORDER BY tanggal DESC, saldo DESC"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return laporans, err
	}

	for rows.Next() {
		laporan := model.LaporanKeuangan{}
		var tanggal interface{} // Simpan tanggal sebagai interface{} untuk debugging

		// Scan ke variabel interface{}
		err := rows.Scan(&laporan.Id, &tanggal, &laporan.Keterangan, &laporan.Pemasukan, &laporan.Pengeluaran, &laporan.Saldo)
		if err != nil {
			return laporans, err
		}

		// Konversi manual jika diperlukan
		switch v := tanggal.(type) {
		case time.Time:
			laporan.Tanggal = v
		case []byte:
			// Jika tanggal adalah []byte, parse ke time.Time
			tanggalStr := string(v)
			parsedTime, err := time.Parse(time.RFC3339, tanggalStr)
			if err != nil {
				return laporans, fmt.Errorf("failed to parse tanggal: %v", err)
			}
			laporan.Tanggal = parsedTime
		default:
			return laporans, fmt.Errorf("unsupported type for tanggal: %T", v)
		}

		laporans = append(laporans, laporan)
	}
	return laporans, nil
}

// GetLastBalance implements LaporanKeuanganRepo by calculating Total Income - Total Expenditure.
func (l *laporanKeuanganRepoImpl) GetLastBalance(ctx context.Context, tx *sql.Tx) (int64, error) {
	var totalIncome, totalExpenditure sql.NullInt64
	query := "SELECT COALESCE(SUM(pemasukan), 0), COALESCE(SUM(pengeluaran), 0) FROM laporan_keuangan"
	err := tx.QueryRowContext(ctx, query).Scan(&totalIncome, &totalExpenditure)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch previous saldo: %v", err)
	}
	return totalIncome.Int64 - totalExpenditure.Int64, nil
}

// GetLaporanByDateRange mengambil data laporan keuangan berdasarkan rentang tanggal
func (l *laporanKeuanganRepoImpl) GetLaporanByDateRange(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]model.LaporanKeuangan, error) {
	var laporans []model.LaporanKeuangan

	// Query dengan filter rentang tanggal
	query := `SELECT id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo FROM laporan_keuangan WHERE tanggal BETWEEN ? AND ? ORDER BY tanggal ASC, saldo ASC`

	// Eksekusi query dengan parameter startDate dan endDate
	rows, err := tx.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return laporans, fmt.Errorf("failed to fetch laporan by date range: %v", err)
	}
	defer rows.Close()

	// Iterasi hasil query dan masukkan ke dalam slice laporans
	for rows.Next() {
		laporan := model.LaporanKeuangan{}
		var tanggalStr string // Simpan tanggal sebagai string terlebih dahulu

		// Scan ke variabel sementara (tanggalStr)
		err := rows.Scan(&laporan.Id, &tanggalStr, &laporan.Keterangan, &laporan.Pemasukan, &laporan.Pengeluaran, &laporan.Saldo)
		if err != nil {
			return laporans, fmt.Errorf("failed to scan laporan: %v", err)
		}

		// Parsing string ke time.Time
		laporan.Tanggal, err = time.Parse(time.RFC3339, tanggalStr) // Sesuaikan format dengan data di database
		if err != nil {
			return laporans, fmt.Errorf("failed to parse tanggal: %v", err)
		}

		laporans = append(laporans, laporan)
	}

	// Periksa error setelah iterasi
	if err = rows.Err(); err != nil {
		return laporans, fmt.Errorf("error after iterating rows: %v", err)
	}

	return laporans, nil
}

// GetTotalIncome implements LaporanKeuanganRepo using SQL SUM.
func (l *laporanKeuanganRepoImpl) GetTotalIncome(ctx context.Context, tx *sql.Tx) (uint64, error) {
	var total sql.NullInt64
	query := "SELECT COALESCE(SUM(pemasukan), 0) FROM laporan_keuangan"
	err := tx.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch total income: %v", err)
	}
	if total.Valid && total.Int64 > 0 {
		return uint64(total.Int64), nil
	}
	return 0, nil
}

// GetTotalExpenditure implements LaporanKeuanganRepo using SQL SUM.
func (l *laporanKeuanganRepoImpl) GetTotalExpenditure(ctx context.Context, tx *sql.Tx) (uint64, error) {
	var total sql.NullInt64
	query := "SELECT COALESCE(SUM(pengeluaran), 0) FROM laporan_keuangan"
	err := tx.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch total expenditure: %v", err)
	}
	if total.Valid && total.Int64 > 0 {
		return uint64(total.Int64), nil
	}
	return 0, nil
}
