package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type AdminRepo interface {
	SignUp(ctx context.Context, tx *sql.Tx, admin model.Admin) (model.Admin, error)
	FindByNik(ctx context.Context, tx *sql.Tx, nik string) (model.Admin, error)
}
