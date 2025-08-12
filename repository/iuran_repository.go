package repository

import (
	"database/sql"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

type IuranRepository interface {
	Create(iuran dto.IuranRequest) (error)
	Update(id string, iuran dto.IuranRequest) error
	GetAll() ([]dto.IuranResponse, error)
	GetByID(id string) (dto.IuranResponse, error)
	Delete(id int64) error
}

type iuranRepository struct {
	db *sql.DB
}

func NewIuranRepo(db *sql.DB) IuranRepository {
	return &iuranRepository{db: db}
}

func (r *iuranRepository) Create(iuran dto.IuranRequest) (error) {
	query := `INSERT INTO iuran (id, nama, jumlah, tanggal, keterangan) VALUES (uuid(), ?, ?, ?, ?)`
	_, err := r.db.Exec(query, iuran.Nama, iuran.Jumlah, iuran.Tanggal, iuran.Keterangan)
	if err != nil {
		return err
	}
	return nil
}

func (r *iuranRepository) Update(id string, iuran dto.IuranRequest) error {
	query := `UPDATE iuran SET nama=?, jumlah=?, tanggal=?, keterangan=? WHERE id=?`
	_, err := r.db.Exec(query, iuran.Nama, iuran.Jumlah, iuran.Tanggal, iuran.Keterangan, id)
	return err
}

func (r *iuranRepository) GetAll() ([]dto.IuranResponse, error) {
	query := `SELECT id, nama, jumlah, tanggal, keterangan FROM iuran`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var iurans []dto.IuranResponse
	for rows.Next() {
		var i dto.IuranResponse
		err := rows.Scan(&i.ID, &i.Nama, &i.Jumlah, &i.Tanggal, &i.Keterangan)
		if err != nil {
			return nil, err
		}
		iurans = append(iurans, i)
	}
	return iurans, nil
}

func (r *iuranRepository) GetByID(id string) (dto.IuranResponse, error) {
	query := `SELECT id, nama, jumlah, tanggal, keterangan FROM iuran WHERE id=?`
	var i dto.IuranResponse
	err := r.db.QueryRow(query, id).Scan(&i.ID, &i.Nama, &i.Jumlah, &i.Tanggal, &i.Keterangan)
	if err != nil {
		return i, err
	}
	return i, nil
}

func (r *iuranRepository) Delete(id int64) error {
	query := `DELETE FROM iuran WHERE id=?`
	_, err := r.db.Exec(query, id)
	return err
}
