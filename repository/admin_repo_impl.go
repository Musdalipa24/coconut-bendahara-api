package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type AdminRepo interface {
	SignUp(ctx context.Context, tx *sql.Tx, admin model.Admin) (model.Admin, error)
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (model.Admin, error)
	UpdateAdmin(ctx context.Context, tx *sql.Tx, admin model.Admin) (model.Admin, error)
}

type adminRepoImpl struct {
}

func NewAdminRepo() AdminRepo {
	return &adminRepoImpl{}
}

// SignUp implements AdminRepo.
func (a adminRepoImpl) SignUp(ctx context.Context, tx *sql.Tx, admin model.Admin) (model.Admin, error) {
	query := "INSERT INTO Admin (id, username, password) VALUES (?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, admin.Id, admin.Username, admin.Password)
	if err != nil {
		return admin, err
	}

	return admin, nil
}

// FindById implements AdminRepo.
func (a adminRepoImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (model.Admin, error) {
	query := "SELECT id, username, password FROM Admin WHERE username = ?"

	rows, err := tx.QueryContext(ctx, query, username)
	if err != nil {
		return model.Admin{}, err
	}

	defer rows.Close()
	admin := model.Admin{}
	if rows.Next() {
		err := rows.Scan(&admin.Id, &admin.Username, &admin.Password)
		if err != nil {
			return model.Admin{}, err
		}
		return admin, nil
	} else {
		return admin, errors.New("username not found")
	}
}

// UpdateAdmin implements AdminRepo.
func (a *adminRepoImpl) UpdateAdmin(ctx context.Context, tx *sql.Tx, admin model.Admin) (model.Admin, error) {
	query := "UPDATE Admin SET password = ? WHERE username = ?"

	_, err := tx.ExecContext(ctx, query, admin.Password, admin.Username)
	if err != nil {
		return model.Admin{}, err
	}

	return admin, nil
}