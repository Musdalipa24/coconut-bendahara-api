package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type IuranRepository interface {
	AddMember(ctx context.Context, tx *sql.Tx, member model.Member) (model.Member, error)
	AddIuran(ctx context.Context, tx *sql.Tx, iuran model.Iuran) (model.Iuran, error)
	AddPembayaran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran, member model.Member) (model.PembayaranIuran, error)
	GetPembayaranById(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran, id string) (model.PembayaranIuran, error)
	GetAllMembers(ctx context.Context, tx *sql.Tx) ([]model.Member, error)
	GetMemberById(ctx context.Context, tx *sql.Tx, id_member string) (model.Member, error)
	GetAllIuran(ctx context.Context, tx *sql.Tx) ([]model.PembayaranIuran, error)
	GetIuranByMemberID(ctx context.Context, tx *sql.Tx, memberID string) ([]model.PembayaranIuran, error)
	GetIuranByPeriod(ctx context.Context, tx *sql.Tx, periode string, minggu_ke int, id_member string) ([]model.PembayaranIuran, error)
	GetIuranByPeriodOnly(ctx context.Context, tx *sql.Tx, periode string, mingguKe int) (model.Iuran, error)
	UpdateStatusIuran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran, member model.Member) (model.PembayaranIuran, error)
	DeleteMember(ctx context.Context, tx *sql.Tx, id_member string) error
}

type iuranRepositoryImpl struct {
}

func NewIuranRepository() IuranRepository {
	return &iuranRepositoryImpl{}
}

// AddMember implements IuranRepository.
func (r *iuranRepositoryImpl) AddMember(ctx context.Context, tx *sql.Tx, member model.Member) (model.Member, error) {
	query := "INSERT INTO member (id_member, nra, nama, status) VALUES (?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, member.IdMember, member.NRA, member.Nama, member.Status)
	if err != nil {
		return model.Member{}, err
	}

	return member, nil
}

// AddIuran implements IuranRepository.
func (r *iuranRepositoryImpl) AddPembayaran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran, member model.Member) (model.PembayaranIuran, error) {
	query := "INSERT INTO pembayaran_iuran (id_pembayaran, id_member, id_pemasukan, status, jumlah_bayar, tanggal_bayar, id_iuran) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, pembayaran.IdPembayaran, pembayaran.IdMember, pembayaran.IdPemasukan, pembayaran.Status, pembayaran.JumlahBayar, pembayaran.TanggalBayar, pembayaran.Iuran.IdIuran)
	if err != nil {
		return model.PembayaranIuran{}, err
	}

	keterangan := fmt.Sprintf("Pembayaran Iuran periode %s - minggu ke %d dari %s (%s)", pembayaran.Iuran.Periode.String, pembayaran.Iuran.MingguKe.Int64, member.Nama, member.NRA)

	id_transaksi := uuid.New().String()

	queryTransaksi := `INSERT INTO history_transaksi (id_transaksi, tanggal, keterangan, jenis_transaksi, nominal) VALUES (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, queryTransaksi, id_transaksi, pembayaran.TanggalBayar.Time, keterangan, "Pemasukan", pembayaran.JumlahBayar.Int64)
	if err != nil {
		return model.PembayaranIuran{}, err
	}

	queryPemasukan := "INSERT INTO pemasukan (id_pemasukan, tanggal, kategori, keterangan, nominal, nota, id_transaksi) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err = tx.ExecContext(ctx, queryPemasukan, pembayaran.IdPemasukan.String, pembayaran.TanggalBayar.Time, "Iuran", keterangan, pembayaran.JumlahBayar.Int64, "no data", id_transaksi)
	if err != nil {
		return model.PembayaranIuran{}, err
	}

	// Ambil saldo terakhir sebelum tanggal pemasukan
	var saldoSebelumnya uint64
	querySaldo := `
		SELECT saldo FROM laporan_keuangan 
		WHERE tanggal <= ?
		ORDER BY tanggal DESC
		LIMIT 1
	`
	err = tx.QueryRowContext(ctx, querySaldo, pembayaran.TanggalBayar).Scan(&saldoSebelumnya)
	if err != nil && err != sql.ErrNoRows {
		return model.PembayaranIuran{}, fmt.Errorf("failed to fetch previous saldo: %v", err)
	}

	// Hitung saldo baru
	saldoBaru := saldoSebelumnya + uint64(pembayaran.JumlahBayar.Int64)

	// Insert laporan keuangan baru
	idLaporan := uuid.New().String()
	queryLaporan := `INSERT INTO laporan_keuangan (id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo, id_transaksi) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, queryLaporan, idLaporan, pembayaran.TanggalBayar, keterangan, pembayaran.JumlahBayar.Int64, 0, saldoBaru, id_transaksi)
	if err != nil {
		return model.PembayaranIuran{}, fmt.Errorf("failed to insert into laporan_keuangan: %v", err)
	}

	// Update saldo semua entri setelah tanggal pemasukan
	queryUpdate := `
		UPDATE laporan_keuangan
		SET saldo = saldo + ?
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate, pembayaran.JumlahBayar, pembayaran.TanggalBayar)
	if err != nil {
		return model.PembayaranIuran{}, fmt.Errorf("failed to update future saldo: %v", err)
	}

	pembayaran.IdTransaksi = sql.NullString{String: id_transaksi, Valid: true}

	return pembayaran, nil
}

// GetPembayaranById implements IuranRepository.
func (r *iuranRepositoryImpl) GetPembayaranById(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran, id_member string) (model.PembayaranIuran, error) {
	query := "SELECT id_pembayaran, id_member, id_iuran, id_pemasukan, status, jumlah_bayar, tanggal_bayar FROM pembayaran_iuran WHERE id_member = ?"

	row := tx.QueryRowContext(ctx, query, id_member)

	err := row.Scan(&pembayaran.IdPembayaran, &pembayaran.IdMember, &pembayaran.Iuran.IdIuran, &pembayaran.IdPemasukan, &pembayaran.Status, &pembayaran.JumlahBayar, &pembayaran.TanggalBayar)
	if err != nil {
		return model.PembayaranIuran{}, err
	}

	return pembayaran, nil
}

// AddIuran implements IuranRepository.
func (r *iuranRepositoryImpl) AddIuran(ctx context.Context, tx *sql.Tx, iuran model.Iuran) (model.Iuran, error) {
	query := "INSERT INTO iuran (id_iuran, periode, minggu_ke) VALUES (?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, iuran.IdIuran, iuran.Periode, iuran.MingguKe.Int64)
	if err != nil {
		return model.Iuran{}, err
	}

	return iuran, nil
}

// GetAllMembers implements MemberRepo.
func (r *iuranRepositoryImpl) GetAllMembers(ctx context.Context, tx *sql.Tx) ([]model.Member, error) {
	query := `
		SELECT 
			id_member,
			nra,
			nama,
			status,
			created_at,
			updated_at
		FROM member ORDER BY id_member ASC
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []model.Member{}, err
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var member model.Member
		err := rows.Scan(
			&member.IdMember,
			&member.NRA,
			&member.Nama,
			&member.Status,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return []model.Member{}, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *iuranRepositoryImpl) GetAllIuran(ctx context.Context, tx *sql.Tx) ([]model.PembayaranIuran, error) {
	query := `
		SELECT 
			pi.id_pembayaran,
			pi.id_member,
			pi.id_pemasukan,
			pi.status,
			pi.jumlah_bayar,
			pi.tanggal_bayar,
			i.id_iuran,
			i.periode,
			i.minggu_ke
		FROM pembayaran_iuran pi
		LEFT JOIN iuran i ON pi.id_iuran = i.id_iuran
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []model.PembayaranIuran{}, err
	}
	defer rows.Close()

	var iurans []model.PembayaranIuran
	for rows.Next() {
		var pembayaran model.PembayaranIuran
		var iuran model.Iuran

		var idPembayaran, idMember, idPemasukan, statusPembayaran sql.NullString
		var tanggalBayar sql.NullTime
		var idIuran, periode sql.NullString
		var mingguKe, jumlah_bayar sql.NullInt64

		err := rows.Scan(
			&idPembayaran,
			&idMember,
			&idPemasukan,
			&statusPembayaran,
			&jumlah_bayar,
			&tanggalBayar,
			&idIuran,
			&periode,
			&mingguKe,
		)
		if err != nil {
			return iurans, err
		}

		pembayaran.IdPembayaran = idPembayaran
		pembayaran.IdMember = idMember
		pembayaran.Status = statusPembayaran
		if tanggalBayar.Valid {
			pembayaran.TanggalBayar = tanggalBayar
		}

		iuran.IdIuran = idIuran
		iuran.Periode = periode
		if mingguKe.Valid {
			iuran.MingguKe = mingguKe
		}

		pembayaran.Iuran = iuran
		iurans = append(iurans, pembayaran)
	}

	return iurans, nil
}

func (i *iuranRepositoryImpl) GetIuranByMemberID(ctx context.Context, tx *sql.Tx, memberID string) ([]model.PembayaranIuran, error) {
	query := `
		SELECT 
			pi.id_pembayaran,
			pi.status,
			pi.jumlah_bayar,
			pi.tanggal_bayar,
			i.id_iuran,
			i.periode,
			i.minggu_ke
		FROM pembayaran_iuran pi
		LEFT JOIN iuran i ON pi.id_iuran = i.id_iuran
		WHERE pi.id_member = ?
	`

	rows, err := tx.QueryContext(ctx, query, memberID)
	if err != nil {
		return []model.PembayaranIuran{}, err
	}
	defer rows.Close()

	var result []model.PembayaranIuran
	for rows.Next() {
		var pembayaran model.PembayaranIuran
		var iuran model.Iuran

		var idPembayaran, statusPembayaran sql.NullString
		var tanggalBayar sql.NullTime
		var idIuran, periode sql.NullString
		var mingguKe, jumlah_bayar sql.NullInt64

		err := rows.Scan(
			&idPembayaran,
			&statusPembayaran,
			&jumlah_bayar,
			&tanggalBayar,
			&idIuran,
			&periode,
			&mingguKe,
		)
		if err != nil {
			return result, err
		}

		if idPembayaran.Valid {
			pembayaran.IdPembayaran = idPembayaran
		}
		if statusPembayaran.Valid {
			pembayaran.Status = statusPembayaran
		}
		if tanggalBayar.Valid {
			pembayaran.TanggalBayar = tanggalBayar
		}
		if jumlah_bayar.Valid {
			pembayaran.JumlahBayar = jumlah_bayar
		}

		if idIuran.Valid {
			iuran.IdIuran = idIuran
		}
		if periode.Valid {
			iuran.Periode = periode
		}
		if mingguKe.Valid {
			iuran.MingguKe = mingguKe
		}

		pembayaran.Iuran = iuran

		result = append(result, pembayaran)
	}

	return result, nil
}

// GetMember implements IuranRepository.
func (i *iuranRepositoryImpl) GetMemberById(ctx context.Context, tx *sql.Tx, id_member string) (model.Member, error) {
	query := `
		SELECT 
			m.id_member,
			m.nra,
			m.nama,
			m.status,
			m.created_at,
			m.updated_at,
			pi.id_pembayaran,
			pi.status,
			pi.jumlah_bayar,
			pi.tanggal_bayar,
			i.id_iuran,
			i.periode,
			i.minggu_ke
		FROM member m
		LEFT JOIN pembayaran_iuran pi ON m.id_member = pi.id_member
		LEFT JOIN iuran i ON pi.id_iuran = i.id_iuran
		WHERE m.id_member = ?
	`

	rows, err := tx.QueryContext(ctx, query, id_member)
	if err != nil {
		return model.Member{}, err
	}
	defer rows.Close()

	var member model.Member
	for rows.Next() {
		var iuran model.Iuran
		var pembayaran model.PembayaranIuran

		var idPembayaran, statusPembayaran sql.NullString
		var tanggalBayar sql.NullTime
		var idIuran, periode sql.NullString
		var mingguKe, jumlah_bayar sql.NullInt64

		err := rows.Scan(
			&member.IdMember,
			&member.NRA,
			&member.Nama,
			&member.Status,
			&member.CreatedAt,
			&member.UpdatedAt,
			&idPembayaran,
			&statusPembayaran,
			&jumlah_bayar,
			&tanggalBayar,
			&idIuran,
			&periode,
			&mingguKe,
		)
		if err != nil {
			return model.Member{}, err
		}

		pembayaran.IdPembayaran = idPembayaran
		pembayaran.Status = statusPembayaran
		if tanggalBayar.Valid {
			pembayaran.TanggalBayar = tanggalBayar
		}
		if jumlah_bayar.Valid {
			pembayaran.JumlahBayar = jumlah_bayar
		}

		iuran.IdIuran = idIuran
		iuran.Periode = periode
		iuran.MingguKe = mingguKe

		pembayaran.Iuran = iuran
		member.PembayaranIurans = append(member.PembayaranIurans, pembayaran)
	}

	return member, nil
}

// GetIuranByPeriod implements IuranRepository.
func (r *iuranRepositoryImpl) GetIuranByPeriod(ctx context.Context, tx *sql.Tx, periode string, minggu_ke int, id_member string) ([]model.PembayaranIuran, error) {
	query := `
		SELECT 
    		pi.id_pembayaran,
    		pi.status,
			pi.jumlah_bayar,
    		pi.tanggal_bayar,
    		pi.id_member,
    		i.id_iuran,
    		i.periode,
    		i.minggu_ke
		FROM pembayaran_iuran pi
		LEFT JOIN iuran i 
    		ON pi.id_iuran = i.id_iuran
		WHERE (i.periode, i.minggu_ke, pi.id_member) IN (
    		(?, ?, ?)
		)
		ORDER BY pi.id_pembayaran;
	`

	rows, err := tx.QueryContext(ctx, query, periode, minggu_ke, id_member)
	if err != nil {
		return []model.PembayaranIuran{}, err
	}
	defer rows.Close()

	var iurans []model.PembayaranIuran
	for rows.Next() {
		var pembayaran model.PembayaranIuran
		var iuran model.Iuran
		var idPembayaran, statusPembayaran sql.NullString
		var tanggalBayar sql.NullTime
		var idMember, idIuran, periode sql.NullString
		var mingguKe, jumlah_bayar sql.NullInt64
		err := rows.Scan(
			&idPembayaran,
			&statusPembayaran,
			&jumlah_bayar,
			&tanggalBayar,
			&idMember,
			&idIuran,
			&periode,
			&mingguKe,
		)
		if err != nil {
			return iurans, err
		}
		pembayaran.IdPembayaran = idPembayaran
		pembayaran.Status = statusPembayaran
		pembayaran.JumlahBayar = jumlah_bayar
		pembayaran.TanggalBayar = tanggalBayar
		pembayaran.IdMember = idMember
		iuran.IdIuran = idIuran
		iuran.Periode = periode
		iuran.MingguKe = mingguKe
		pembayaran.Iuran = iuran
		iurans = append(iurans, pembayaran)
	}
	return iurans, nil
}

// GetIuranByPeriodOnly mengambil data iuran berdasarkan periode & minggu_ke saja
func (r *iuranRepositoryImpl) GetIuranByPeriodOnly(ctx context.Context, tx *sql.Tx, periode string, mingguKe int) (model.Iuran, error) {
	query := `
		SELECT id_iuran, periode, minggu_ke
		FROM iuran
		WHERE periode = ? AND minggu_ke = ?
		LIMIT 1
	`

	row := tx.QueryRowContext(ctx, query, periode, mingguKe)

	var iuran model.Iuran
	err := row.Scan(
		&iuran.IdIuran,
		&iuran.Periode,
		&iuran.MingguKe,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Iuran{}, nil
		}
		return model.Iuran{}, err
	}

	return iuran, nil
}

// UpdateIuran implements IuranRepository.
func (r *iuranRepositoryImpl) UpdateStatusIuran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran, member model.Member) (model.PembayaranIuran, error) {
	// PERBAIKAN: Update pembayaran_iuran langsung tanpa JOIN yang tidak perlu
	query := "UPDATE pembayaran_iuran SET status = ?, tanggal_bayar = ?, jumlah_bayar = ? WHERE id_member = ? AND id_iuran = ?"

	_, err := tx.ExecContext(ctx, query, pembayaran.Status, pembayaran.TanggalBayar, pembayaran.JumlahBayar, pembayaran.IdMember, pembayaran.Iuran.IdIuran)
	if err != nil {
		return model.PembayaranIuran{}, err
	}

	var oldNominal uint64
	var tanggalRaw []byte
	var idTransaksi string
	queryFetch := `
		SELECT nominal, tanggal, id_transaksi 
		FROM pemasukan 
		WHERE id_pemasukan = ?
	`
	err = tx.QueryRowContext(ctx, queryFetch, pembayaran.IdPemasukan).Scan(&oldNominal, &tanggalRaw, &idTransaksi)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.PembayaranIuran{}, fmt.Errorf("pemasukan with id %s not found", pembayaran.IdPemasukan.String)
		}
		return model.PembayaranIuran{}, fmt.Errorf("failed to fetch previous pemasukan: %v", err)
	}

	keterangan := fmt.Sprintf("Pembayaran Iuran periode %s - minggu ke %d dari %s (%s)", pembayaran.Iuran.Periode.String, pembayaran.Iuran.MingguKe.Int64, member.Nama, member.NRA)

	tanggalStr := string(tanggalRaw)
	oldTanggal, err := time.Parse(time.RFC3339, tanggalStr)
	if err != nil {
		return model.PembayaranIuran{}, fmt.Errorf("failed to parse old tanggal: %v", err)
	}

	// Hitung selisih nominal
	nominalDiff := int64(pembayaran.JumlahBayar.Int64) - int64(oldNominal)

	// Perbarui tabel pemasukan
	queryPemasukan := `
		UPDATE pemasukan 
		SET tanggal = ?, kategori = ?, keterangan = ?, nominal = ?, nota = ? 
		WHERE id_pemasukan = ?
	`
	_, err = tx.ExecContext(ctx, queryPemasukan, pembayaran.TanggalBayar, "Iuran", keterangan, pembayaran.JumlahBayar, "no data", pembayaran.IdPemasukan)
	if err != nil {
		return model.PembayaranIuran{}, fmt.Errorf("failed to update pemasukan: %v", err)
	}

	// Perbarui tabel history_transaksi
	queryHistory := `
		UPDATE history_transaksi 
		SET tanggal = ?, keterangan = ?, nominal = ? 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, pembayaran.TanggalBayar, keterangan, pembayaran.JumlahBayar, idTransaksi)
	if err != nil {
		return model.PembayaranIuran{}, fmt.Errorf("failed to update history_transaksi: %v", err)
	}

	// PERBAIKAN: Update laporan_keuangan dengan benar
	// Jika tanggal tidak berubah, hanya update nominal di record yang sama
	if oldTanggal.Equal(pembayaran.TanggalBayar.Time) {
		// Update record laporan_keuangan yang sama
		queryLaporan := `
			UPDATE laporan_keuangan 
			SET keterangan = ?, pemasukan = ?, saldo = saldo + ? 
			WHERE id_transaksi = ?
		`
		_, err = tx.ExecContext(ctx, queryLaporan, keterangan, pembayaran.JumlahBayar, nominalDiff, idTransaksi)
		if err != nil {
			return model.PembayaranIuran{}, fmt.Errorf("failed to update laporan_keuangan: %v", err)
		}

		// Update saldo untuk semua record setelah tanggal ini
		if nominalDiff != 0 {
			queryUpdateFuture := `
				UPDATE laporan_keuangan 
				SET saldo = saldo + ? 
				WHERE tanggal > ?
			`
			_, err = tx.ExecContext(ctx, queryUpdateFuture, nominalDiff, pembayaran.TanggalBayar)
			if err != nil {
				return model.PembayaranIuran{}, fmt.Errorf("failed to update future laporan_keuangan: %v", err)
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
			return model.PembayaranIuran{}, fmt.Errorf("failed to delete old laporan_keuangan: %v", err)
		}

		// 2. Update saldo untuk record setelah tanggal lama (kurangi nominal lama)
		queryUpdateAfterOld := `
			UPDATE laporan_keuangan 
			SET saldo = saldo - ? 
			WHERE tanggal > ?
		`
		_, err = tx.ExecContext(ctx, queryUpdateAfterOld, oldNominal, oldTanggal)
		if err != nil {
			return model.PembayaranIuran{}, fmt.Errorf("failed to update records after old date: %v", err)
		}

		// 3. Ambil saldo terakhir sebelum tanggal baru
		var saldoSebelumnya uint64
		querySaldo := `
			SELECT saldo FROM laporan_keuangan 
			WHERE tanggal <= ?
			ORDER BY tanggal DESC
			LIMIT 1
		`
		err = tx.QueryRowContext(ctx, querySaldo, pembayaran.TanggalBayar).Scan(&saldoSebelumnya)
		if err != nil && err != sql.ErrNoRows {
			return model.PembayaranIuran{}, fmt.Errorf("failed to fetch previous saldo: %v", err)
		}

		// 4. Hitung saldo baru
		saldoBaru := saldoSebelumnya + uint64(pembayaran.JumlahBayar.Int64)

		// 5. Insert record laporan_keuangan baru
		idLaporan := uuid.New().String()
		queryInsertNew := `
			INSERT INTO laporan_keuangan 
			(id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo, id_transaksi)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.ExecContext(ctx, queryInsertNew, idLaporan, pembayaran.TanggalBayar, keterangan, pembayaran.JumlahBayar, 0, saldoBaru, idTransaksi)
		if err != nil {
			return model.PembayaranIuran{}, fmt.Errorf("failed to insert new laporan_keuangan: %v", err)
		}

		// 6. Update saldo untuk record setelah tanggal baru (tambah nominal baru)
		queryUpdateAfterNew := `
			UPDATE laporan_keuangan 
			SET saldo = saldo + ? 
			WHERE tanggal > ?
		`
		_, err = tx.ExecContext(ctx, queryUpdateAfterNew, pembayaran.JumlahBayar, pembayaran.TanggalBayar)
		if err != nil {
			return model.PembayaranIuran{}, fmt.Errorf("failed to update records after new date: %v", err)
		}
	}

	return pembayaran, nil
}

// DeleteMember implements IuranRepository.
func (r *iuranRepositoryImpl) DeleteMember(ctx context.Context, tx *sql.Tx, id_member string) error {
	// PERBAIKAN: Hapus semua data terkait dengan member

	// 1. Ambil semua pembayaran iuran dari member ini untuk rollback laporan keuangan
	getPembayaran, err := r.GetIuranByMemberID(ctx, tx, id_member)
	if err != nil {
		return fmt.Errorf("failed to get pembayaran iuran for member: %v", err)
	}

	// 2. Untuk setiap pembayaran, rollback laporan keuangan
	for _, pembayaran := range getPembayaran {
		if pembayaran.IdPemasukan.Valid {
			// Ambil data pemasukan untuk mendapatkan nominal dan tanggal
			var nominal uint64
			var tanggalRaw []byte
			var idTransaksi string
			queryFetch := `
				SELECT nominal, tanggal, id_transaksi 
				FROM pemasukan 
				WHERE id_pemasukan = ?
			`
			err = tx.QueryRowContext(ctx, queryFetch, pembayaran.IdPemasukan.String).Scan(&nominal, &tanggalRaw, &idTransaksi)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("failed to fetch pemasukan data: %v", err)
			}

			if err != sql.ErrNoRows {
				// Parse tanggal
				tanggalStr := string(tanggalRaw)
				tanggalTime, err := time.Parse(time.RFC3339, tanggalStr)
				if err != nil {
					return fmt.Errorf("failed to parse tanggal: %v", err)
				}

				// Hapus dari laporan_keuangan
				queryDeleteLaporan := `
					DELETE FROM laporan_keuangan 
					WHERE id_transaksi = ?
				`
				_, err = tx.ExecContext(ctx, queryDeleteLaporan, idTransaksi)
				if err != nil {
					return fmt.Errorf("failed to delete from laporan_keuangan: %v", err)
				}

				// Update saldo untuk semua record setelah tanggal ini (kurangi nominal)
				queryUpdateSaldo := `
					UPDATE laporan_keuangan
					SET saldo = saldo - ?
					WHERE tanggal > ?
				`
				_, err = tx.ExecContext(ctx, queryUpdateSaldo, nominal, tanggalTime)
				if err != nil {
					return fmt.Errorf("failed to update future saldo: %v", err)
				}

				// Hapus dari history_transaksi
				queryDeleteHistory := `
					DELETE FROM history_transaksi 
					WHERE id_transaksi = ?
				`
				_, err = tx.ExecContext(ctx, queryDeleteHistory, idTransaksi)
				if err != nil {
					return fmt.Errorf("failed to delete from history_transaksi: %v", err)
				}

				// Hapus dari pemasukan
				queryDeletePemasukan := `
					DELETE FROM pemasukan 
					WHERE id_pemasukan = ?
				`
				_, err = tx.ExecContext(ctx, queryDeletePemasukan, pembayaran.IdPemasukan.String)
				if err != nil {
					return fmt.Errorf("failed to delete from pemasukan: %v", err)
				}
			}
		}
	}

	// 3. Hapus semua pembayaran_iuran dari member ini
	queryDeletePembayaran := "DELETE FROM pembayaran_iuran WHERE id_member = ?"
	_, err = tx.ExecContext(ctx, queryDeletePembayaran, id_member)
	if err != nil {
		return fmt.Errorf("failed to delete pembayaran_iuran: %v", err)
	}

	// 4. Terakhir hapus member
	queryDeleteMember := "DELETE FROM member WHERE id_member = ?"
	_, err = tx.ExecContext(ctx, queryDeleteMember, id_member)
	if err != nil {
		return fmt.Errorf("failed to delete member: %v", err)
	}

	return nil
}
