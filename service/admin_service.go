package service

import (
	"context"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

type AdminService interface {
	SignUp(ctx context.Context, adminRequest dto.AdminRequest) (dto.AdminResponse, error)
	SignIn(ctx context.Context, loginRequest dto.LoginRequest) (string, error)
	GetAdminByNik(ctx context.Context, nik string) (dto.AdminResponse, error)
	GenerateJWT(email string) (string, error)
}