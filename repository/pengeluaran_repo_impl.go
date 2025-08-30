package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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

type pengeluaranRepoImpl struct {
}

func NewPengeluaranRepo() PengeluaranRepo {
	return &pengeluaranRepoImpl{}
}

// AddPengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) AddPengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran) (model.Pengeluaran, error) {
	// Validasi tanggal
	if pengeluaran.Tanggal.IsZero() {
		return pengeluaran, fmt.Errorf("tanggal cannot be zero")
	}

	// PERBAIKAN: Ambil saldo terakhir sebelum tanggal pengeluaran DI AWAL
	var saldoSebelumnya uint64
	querySaldo := `
		SELECT saldo FROM laporan_keuangan 
		WHERE tanggal <= ?
		ORDER BY tanggal DESC
		LIMIT 1
	`
	err := tx.QueryRowContext(ctx, querySaldo, pengeluaran.Tanggal).Scan(&saldoSebelumnya)
	if err != nil && err != sql.ErrNoRows {
		return pengeluaran, fmt.Errorf("failed to fetch previous saldo: %v", err)
	}

	// PERBAIKAN: Validasi saldo cukup untuk pengeluaran SEBELUM insert apapun
	if pengeluaran.Nominal > saldoSebelumnya {
		return pengeluaran, fmt.Errorf("saldo tidak mencukupi. Saldo tersedia: %d, diperlukan: %d", saldoSebelumnya, pengeluaran.Nominal)
	}

	idTransaksi := uuid.New().String()

	// Insert ke history_transaksi
	queryTransaksi := `
		INSERT INTO history_transaksi (id_transaksi, tanggal, keterangan, jenis_transaksi, nominal)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryTransaksi, idTransaksi, pengeluaran.Tanggal, pengeluaran.Keterangan, "Pengeluaran", pengeluaran.Nominal)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to insert into history_transaksi: %v", err)
	}

	// Insert ke tabel pengeluaran
	queryPengeluaran := `
		INSERT INTO pengeluaran (id_pengeluaran, tanggal, nota, nominal, keterangan, id_transaksi)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryPengeluaran, pengeluaran.Id, pengeluaran.Tanggal, pengeluaran.Nota, pengeluaran.Nominal, pengeluaran.Keterangan, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to insert into pengeluaran: %v", err)
	}

	// Hitung saldo baru
	saldoBaru := saldoSebelumnya - pengeluaran.Nominal

	// Insert laporan keuangan baru
	idLaporan := uuid.New().String()
	queryLaporan := `
		INSERT INTO laporan_keuangan 
		(id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo, id_transaksi)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryLaporan, idLaporan, pengeluaran.Tanggal, pengeluaran.Keterangan, 0, pengeluaran.Nominal, saldoBaru, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to insert into laporan_keuangan: %v", err)
	}

	// Update saldo semua entri setelah tanggal pengeluaran
	queryUpdate := `
		UPDATE laporan_keuangan
		SET saldo = saldo - ?
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate, pengeluaran.Nominal, pengeluaran.Tanggal)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to update future saldo: %v", err)
	}

	pengeluaran.IdTransaksi = idTransaksi
	return pengeluaran, nil
}

// UpdatePengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) UpdatePengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran, id string) (model.Pengeluaran, error) {
	// Pastikan tanggal sudah dalam format time.Time
	if pengeluaran.Tanggal.IsZero() {
		return pengeluaran, fmt.Errorf("tanggal cannot be zero")
	}

	// Ambil data pengeluaran sebelumnya untuk mendapatkan nominal lama, tanggal lama, dan id_transaksi
	var oldNominal uint64
	var tanggalRaw []byte
	var idTransaksi string
	queryFetch := `
		SELECT nominal, tanggal, id_transaksi 
		FROM pengeluaran 
		WHERE id_pengeluaran = ?
	`
	err := tx.QueryRowContext(ctx, queryFetch, id).Scan(&oldNominal, &tanggalRaw, &idTransaksi)
	if err != nil {
		if err == sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("pengeluaran with id %s not found", id)
		}
		return pengeluaran, fmt.Errorf("failed to fetch previous pengeluaran: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	oldTanggal, err := time.Parse(time.RFC3339, tanggalStr)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to parse old tanggal: %v", err)
	}

	// PERBAIKAN: Validasi saldo untuk update
	// Jika nominal baru lebih besar dari nominal lama, cek apakah saldo mencukupi
	if pengeluaran.Nominal > oldNominal {
		// Ambil saldo sebelum pengeluaran lama
		var saldoSebelumnya uint64
		querySaldo := `
			SELECT saldo FROM laporan_keuangan 
			WHERE tanggal <= ? AND id_transaksi != ?
			ORDER BY tanggal DESC
			LIMIT 1
		`
		err = tx.QueryRowContext(ctx, querySaldo, oldTanggal, idTransaksi).Scan(&saldoSebelumnya)
		if err != nil && err != sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("failed to fetch previous saldo: %v", err)
		}

		// Tambahkan kembali nominal lama untuk mendapatkan saldo asli
		saldoAsli := saldoSebelumnya + oldNominal

		// Validasi apakah saldo mencukupi untuk nominal baru
		if pengeluaran.Nominal > saldoAsli {
			return pengeluaran, fmt.Errorf("saldo tidak mencukupi untuk update. Saldo tersedia: %d, diperlukan: %d", saldoAsli, pengeluaran.Nominal)
		}
	}

	// Hitung selisih nominal (positif jika nominal baru lebih kecil, negatif jika lebih besar)
	nominalDiff := int64(oldNominal) - int64(pengeluaran.Nominal)

	// Perbarui tabel pengeluaran
	queryPengeluaran := `
		UPDATE pengeluaran 
		SET tanggal = ?, nota = ?, nominal = ?, keterangan = ? 
		WHERE id_pengeluaran = ?
	`
	_, err = tx.ExecContext(ctx, queryPengeluaran, pengeluaran.Tanggal, pengeluaran.Nota, pengeluaran.Nominal, pengeluaran.Keterangan, id)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to update pengeluaran: %v", err)
	}

	// Perbarui tabel history_transaksi
	queryHistory := `
		UPDATE history_transaksi 
		SET tanggal = ?, keterangan = ?, nominal = ? 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, pengeluaran.Tanggal, pengeluaran.Keterangan, pengeluaran.Nominal, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to update history_transaksi: %v", err)
	}

	// PERBAIKAN: Update laporan_keuangan dengan benar
	// Jika tanggal tidak berubah, hanya update nominal di record yang sama
	if oldTanggal.Equal(pengeluaran.Tanggal) {
		// Update record laporan_keuangan yang sama
		queryLaporan := `
			UPDATE laporan_keuangan 
			SET keterangan = ?, pengeluaran = ?, saldo = saldo + ? 
			WHERE id_transaksi = ?
		`
		_, err = tx.ExecContext(ctx, queryLaporan, pengeluaran.Keterangan, pengeluaran.Nominal, nominalDiff, idTransaksi)
		if err != nil {
			return pengeluaran, fmt.Errorf("failed to update laporan_keuangan: %v", err)
		}

		// Update saldo untuk semua record setelah tanggal ini
		if nominalDiff != 0 {
			queryUpdateFuture := `
				UPDATE laporan_keuangan 
				SET saldo = saldo + ? 
				WHERE tanggal > ?
			`
			_, err = tx.ExecContext(ctx, queryUpdateFuture, nominalDiff, pengeluaran.Tanggal)
			if err != nil {
				return pengeluaran, fmt.Errorf("failed to update future laporan_keuangan: %v", err)
			}
		}
	} else {
		// Jika tanggal berubah, hapus record lama dan buat record baru

		// 1. Hapus record laporan_keuangan lama
		queryDeleteOld := `
			DELETE FROM laporan_keuangan 
			WHERE id_transaksi = ?
		`
		_, err = tx.ExecContext(ctx, queryDeleteOld, idTransaksi)
		if err != nil {
			return pengeluaran, fmt.Errorf("failed to delete old laporan_keuangan: %v", err)
		}

		// 2. Update saldo untuk record setelah tanggal lama (tambah kembali nominal lama)
		queryUpdateAfterOld := `
			UPDATE laporan_keuangan 
			SET saldo = saldo + ? 
			WHERE tanggal > ?
		`
		_, err = tx.ExecContext(ctx, queryUpdateAfterOld, oldNominal, oldTanggal)
		if err != nil {
			return pengeluaran, fmt.Errorf("failed to update records after old date: %v", err)
		}

		// 3. Ambil saldo terakhir sebelum tanggal baru
		var saldoSebelumnya uint64
		querySaldo := `
			SELECT saldo FROM laporan_keuangan 
			WHERE tanggal <= ?
			ORDER BY tanggal DESC
			LIMIT 1
		`
		err = tx.QueryRowContext(ctx, querySaldo, pengeluaran.Tanggal).Scan(&saldoSebelumnya)
		if err != nil && err != sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("failed to fetch previous saldo: %v", err)
		}

		// 4. Hitung saldo baru
		saldoBaru := saldoSebelumnya - pengeluaran.Nominal

		// 5. Insert record laporan_keuangan baru
		idLaporan := uuid.New().String()
		queryInsertNew := `
			INSERT INTO laporan_keuangan 
			(id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo, id_transaksi)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.ExecContext(ctx, queryInsertNew, idLaporan, pengeluaran.Tanggal, pengeluaran.Keterangan, 0, pengeluaran.Nominal, saldoBaru, idTransaksi)
		if err != nil {
			return pengeluaran, fmt.Errorf("failed to insert new laporan_keuangan: %v", err)
		}

		// 6. Update saldo untuk record setelah tanggal baru (kurangi nominal baru)
		queryUpdateAfterNew := `
			UPDATE laporan_keuangan 
			SET saldo = saldo - ? 
			WHERE tanggal > ?
		`
		_, err = tx.ExecContext(ctx, queryUpdateAfterNew, pengeluaran.Nominal, pengeluaran.Tanggal)
		if err != nil {
			return pengeluaran, fmt.Errorf("failed to update records after new date: %v", err)
		}
	}

	return pengeluaran, nil
}

// GetPengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) GetPengeluaran(ctx context.Context, tx *sql.Tx, page int, pageSize int) ([]model.Pengeluaran, int, error) {
	// Hitung offset
	offset := (page - 1) * pageSize

	// Query untuk mendapatkan data dengan pagination
	query := "SELECT id_pengeluaran, tanggal, nota, nominal, keterangan FROM pengeluaran ORDER BY tanggal DESC LIMIT ? OFFSET ?"
	rows, err := tx.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch pengeluaran: %v", err)
	}
	defer rows.Close()

	var pengeluaranSlice []model.Pengeluaran
	for rows.Next() {
		pengeluaran := model.Pengeluaran{}
		var tanggalRaw []byte

		err := rows.Scan(&pengeluaran.Id, &tanggalRaw, &pengeluaran.Nota, &pengeluaran.Nominal, &pengeluaran.Keterangan)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pengeluaran: %v", err)
		}

		// Konversi tanggalRaw ke time.Time
		tanggalStr := string(tanggalRaw)
		parsedTime, err := time.Parse(time.RFC3339, tanggalStr)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse tanggal: %v", err)
		}
		pengeluaran.Tanggal = parsedTime

		pengeluaranSlice = append(pengeluaranSlice, pengeluaran)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error after iterating rows: %v", err)
	}

	// Query untuk mendapatkan total data
	var total int
	countQuery := "SELECT COUNT(*) FROM pengeluaran"
	err = tx.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total pengeluaran: %v", err)
	}

	return pengeluaranSlice, total, nil
}

// FindById implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (model.Pengeluaran, error) {
	query := "SELECT id_pengeluaran, tanggal, nota, nominal, keterangan FROM pengeluaran WHERE id_pengeluaran = ?"
	row := tx.QueryRowContext(ctx, query, id)

	pengeluaran := model.Pengeluaran{}
	var tanggalRaw []byte

	err := row.Scan(&pengeluaran.Id, &tanggalRaw, &pengeluaran.Nota, &pengeluaran.Nominal, &pengeluaran.Keterangan)
	if err != nil {
		if err == sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("pengeluaran not found")
		}
		return pengeluaran, fmt.Errorf("failed to scan pengeluaran: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	parsedTime, err := time.Parse(time.RFC3339, tanggalStr)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to parse tanggal: %v", err)
	}
	pengeluaran.Tanggal = parsedTime
	return pengeluaran, nil
}

// DeletePengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) DeletePengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran) (model.Pengeluaran, error) {
	// Validate input
	if pengeluaran.Id == "" {
		return pengeluaran, fmt.Errorf("id_pengeluaran cannot be empty")
	}

	// Fetch id_transaksi, nominal, and tanggal from pengeluaran
	var idTransaksi string
	var nominal uint64
	var tanggalRaw []byte
	queryFetch := `
		SELECT id_transaksi, nominal, tanggal 
		FROM pengeluaran 
		WHERE id_pengeluaran = ?
	`
	err := tx.QueryRowContext(ctx, queryFetch, pengeluaran.Id).Scan(&idTransaksi, &nominal, &tanggalRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("pengeluaran with id %s not found", pengeluaran.Id)
		}
		return pengeluaran, fmt.Errorf("failed to fetch pengeluaran: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	tanggalTime, err := time.Parse(time.RFC3339, tanggalStr)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to parse tanggal: %v", err)
	}

	// Validate nominal
	if nominal <= 0 {
		return pengeluaran, fmt.Errorf("invalid nominal value: %d", nominal)
	}

	// Delete from laporan_keuangan
	queryLaporan := `
		DELETE FROM laporan_keuangan 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryLaporan, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to delete from laporan_keuangan: %v", err)
	}

	// Delete from history_transaksi
	queryHistory := `
		DELETE FROM history_transaksi 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to delete from history_transaksi: %v", err)
	}

	// Delete from pengeluaran
	queryPengeluaran := `
		DELETE FROM pengeluaran 
		WHERE id_pengeluaran = ?
	`
	_, err = tx.ExecContext(ctx, queryPengeluaran, pengeluaran.Id)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to delete from pengeluaran: %v", err)
	}

	// Update saldo for all records after the deleted pengeluaran's tanggal
	queryUpdateSaldo := `
		UPDATE laporan_keuangan
		SET saldo = saldo + ?
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdateSaldo, nominal, tanggalTime)
	if err != nil {
		return pengeluaran, fmt.Errorf("failed to update future saldo: %v", err)
	}

	return pengeluaran, nil
}

// GetPengeluaranByDateRange implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) GetPengeluaranByDateRange(ctx context.Context, tx *sql.Tx, startDate, endDate string, page int, pageSize int) ([]model.Pengeluaran, int, error) {
	// Hitung offset
	offset := (page - 1) * pageSize

	// Query untuk mendapatkan data dengan pagination dan date range
	query := `
		SELECT id_pengeluaran, tanggal, nota, nominal, keterangan 
		FROM pengeluaran 
		WHERE tanggal BETWEEN ? AND ? 
		ORDER BY tanggal DESC 
		LIMIT ? OFFSET ?
	`
	rows, err := tx.QueryContext(ctx, query, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch pengeluaran by date range: %v", err)
	}
	defer rows.Close()

	var pengeluaranSlice []model.Pengeluaran
	for rows.Next() {
		pengeluaran := model.Pengeluaran{}
		var tanggalRaw []byte

		err := rows.Scan(&pengeluaran.Id, &tanggalRaw, &pengeluaran.Nota, &pengeluaran.Nominal, &pengeluaran.Keterangan)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pengeluaran: %v", err)
		}

		// Konversi tanggalRaw ke time.Time
		tanggalStr := string(tanggalRaw)
		parsedTime, err := time.Parse(time.RFC3339, tanggalStr)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse tanggal: %v", err)
		}
		pengeluaran.Tanggal = parsedTime

		pengeluaranSlice = append(pengeluaranSlice, pengeluaran)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error after iterating rows: %v", err)
	}

	// Query untuk mendapatkan total data dalam rentang tanggal
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM pengeluaran 
		WHERE tanggal BETWEEN ? AND ?
	`
	err = tx.QueryRowContext(ctx, countQuery, startDate, endDate).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total pengeluaran in date range: %v", err)
	}

	return pengeluaranSlice, total, nil
}
