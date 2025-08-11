package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type TransaksiRepo interface {
	GetAllTransaction(ctx context.Context, tx *sql.Tx) ([]model.Transaction, error)
	GetLastTransaction(ctx context.Context, tx *sql.Tx) ([]model.Transaction, error)
}
