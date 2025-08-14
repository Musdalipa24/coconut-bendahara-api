package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type IuranRepository interface {
	AddMember(ctx context.Context, tx *sql.Tx, member model.Member) (model.Member, error)
	AddIuran(ctx context.Context, tx *sql.Tx, iuran model.Iuran) (model.Iuran, error)
	AddPembayaran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran) (model.PembayaranIuran, error)
	GetAllMembers(ctx context.Context, tx *sql.Tx) ([]model.Member, error)
	GetMemberById(ctx context.Context, tx *sql.Tx, id_member string) (model.Member, error)
	GetAllIuran(ctx context.Context, tx *sql.Tx) ([]model.PembayaranIuran, error)
	GetIuranByMemberID(ctx context.Context, tx *sql.Tx, memberID string) ([]model.PembayaranIuran, error)
	GetIuranByPeriod(ctx context.Context, tx *sql.Tx, periode string, minggu_ke int) ([]model.PembayaranIuran, error)
	GetIuranByPeriodOnly(ctx context.Context, tx *sql.Tx, periode string, mingguKe int) (model.Iuran, error)
	UpdateStatusIuran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran) (model.PembayaranIuran, error)
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
func (r *iuranRepositoryImpl) AddPembayaran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran) (model.PembayaranIuran, error) {
	query := "INSERT INTO pembayaran_iuran (id_pembayaran, id_member, status, tanggal_bayar, id_iuran) VALUES (?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, pembayaran.IdPembayaran, pembayaran.IdMember, pembayaran.Status, pembayaran.TanggalBayar, pembayaran.Iuran.IdIuran)
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
		FROM member
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
			pi.status,
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

		var idPembayaran, idMember, statusPembayaran sql.NullString
		var tanggalBayar sql.NullString
		var idIuran, periode sql.NullString
		var mingguKe sql.NullInt64

		err := rows.Scan(
			&idPembayaran,
			&idMember,
			&statusPembayaran,
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

		var idPembayaran, statusPembayaran, tanggalBayar sql.NullString
		var idIuran, periode sql.NullString
		var mingguKe sql.NullInt64

		err := rows.Scan(
			&idPembayaran,
			&statusPembayaran,
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
		var tanggalBayar sql.NullString
		var idIuran, periode sql.NullString
		var mingguKe sql.NullInt64

		err := rows.Scan(
			&member.IdMember,
			&member.NRA,
			&member.Nama,
			&member.Status,
			&member.CreatedAt,
			&member.UpdatedAt,
			&idPembayaran,
			&statusPembayaran,
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

		iuran.IdIuran = idIuran
		iuran.Periode = periode
		iuran.MingguKe = mingguKe

		pembayaran.Iuran = iuran
		member.PembayaranIurans = append(member.PembayaranIurans, pembayaran)
	}

	return member, nil
}

// GetIuranByPeriod implements IuranRepository.
func (r *iuranRepositoryImpl) GetIuranByPeriod(ctx context.Context, tx *sql.Tx, periode string, minggu_ke int) ([]model.PembayaranIuran, error) {
	query := `
		SELECT 
    		pi.id_pembayaran,
    		pi.status,
    		pi.tanggal_bayar,
    		pi.id_member,
    		i.id_iuran,
    		i.periode,
    		i.minggu_ke
		FROM pembayaran_iuran pi
		LEFT JOIN iuran i 
    		ON pi.id_iuran = i.id_iuran
		WHERE (i.periode, i.minggu_ke) IN (
    		(?, ?)
		)
		ORDER BY pi.id_pembayaran;
	`

	rows, err := tx.QueryContext(ctx, query, periode, minggu_ke)
	if err != nil {
		return []model.PembayaranIuran{}, err
	}
	defer rows.Close()

	var iurans []model.PembayaranIuran
	for rows.Next() {
		var pembayaran model.PembayaranIuran
		var iuran model.Iuran
		var idPembayaran, statusPembayaran, tanggalBayar sql.NullString
		var idMember, idIuran, periode sql.NullString
		var mingguKe sql.NullInt64
		err := rows.Scan(
			&idPembayaran,
			&statusPembayaran,
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
			return model.Iuran{}, nil // tidak ada data
		}
		return model.Iuran{}, err
	}

	return iuran, nil
}


// UpdateIuran implements IuranRepository.
func (r *iuranRepositoryImpl) UpdateStatusIuran(ctx context.Context, tx *sql.Tx, pembayaran model.PembayaranIuran) (model.PembayaranIuran, error) {
	query := "UPDATE iuran AS i JOIN pembayaran_iuran AS pi ON i.id_iuran = pi.id_iuran SET pi.status = ?, pi.tanggal_bayar = ? WHERE pi.id_member = ? and i.periode = ? and i.minggu_ke = ?"

	_, err := tx.ExecContext(ctx, query, pembayaran.Status, pembayaran.TanggalBayar, pembayaran.IdMember, pembayaran.Iuran.Periode, pembayaran.Iuran.MingguKe.Int64)
	if err != nil {
		return model.PembayaranIuran{}, err
	}

	return pembayaran, nil
}
