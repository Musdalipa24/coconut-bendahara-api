package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	SignUp(ctx context.Context, adminRequest dto.AdminRequest) (dto.AdminResponse, error)
	SignIn(ctx context.Context, loginRequest dto.LoginRequest) (string, error)
	GetAdminByUsername(ctx context.Context, username string) (dto.AdminResponse, error)
	GenerateJWT(email string) (string, error)
	UpdateAdmin(ctx context.Context, adminRequest dto.UpdateAdminRequest, username string) (dto.AdminResponse, error)
}

type adminServiceImpl struct {
	AdminRepo repository.AdminRepo
	DB        *sql.DB
}

func NewAdminService(adminRepo repository.AdminRepo, db *sql.DB) AdminService {
	return adminServiceImpl{
		AdminRepo: adminRepo,
		DB:        db,
	}
}

// Function for hash password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Function to verify password
func verifyPassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// SignUp implements AdminService.
func (a adminServiceImpl) SignUp(ctx context.Context, adminRequest dto.AdminRequest) (dto.AdminResponse, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Commit()

	hassedPass, err := hashPassword(adminRequest.Password)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to hash password: %v", err)
	}

	admin := model.Admin{
		Id:       uuid.New().String(),
		Username: adminRequest.Username,
		Password: hassedPass,
	}

	createAdmin, err := a.AdminRepo.SignUp(ctx, tx, admin)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to register admin: %v", err)
	}

	return util.ConvertAdminToResponseDTO(createAdmin), nil
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (a adminServiceImpl) GenerateJWT(username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	jwtKey := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtKey))
}

// SignIn implements AdminService.
func (a adminServiceImpl) SignIn(ctx context.Context, loginRequest dto.LoginRequest) (string, error) {
	if loginRequest.Username == "" || loginRequest.Password == "" {
		return "", fmt.Errorf("username or password can't be empty")
	}
	tx, err := a.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Commit()

	admin, err := a.AdminRepo.FindByUsername(ctx, tx, loginRequest.Username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password: %v", err)
	}

	if verifyPassword(admin.Password, loginRequest.Password) {
		fmt.Println("Login berhasil!")
	} else {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := a.GenerateJWT(loginRequest.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

// GetAdminByNik implements AdminService.
func (a adminServiceImpl) GetAdminByUsername(ctx context.Context, username string) (dto.AdminResponse, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	admin, err := a.AdminRepo.FindByUsername(ctx, tx, username)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("admin with username %s not found: %v", username, err)
	}

	return util.ConvertAdminToResponseDTO(admin), nil
}

// UpdateAdmin implements AdminService.
func (a adminServiceImpl) UpdateAdmin(ctx context.Context, adminRequest dto.UpdateAdminRequest, username string) (dto.AdminResponse, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	admin, err := a.AdminRepo.FindByUsername(ctx, tx, username)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("admin with username %s not found: %v", username, err)
	}

	if adminRequest.OldPassword != "" && !verifyPassword(admin.Password, adminRequest.OldPassword) {
		return dto.AdminResponse{}, fmt.Errorf("old password is incorrect")
	}

	admin.Password = adminRequest.Password
	if admin.Password == "" {
		return dto.AdminResponse{}, fmt.Errorf("password can't be empty")
	}
	hassedPass, err := hashPassword(admin.Password)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to hash password: %v", err)
	}
	admin.Password = hassedPass

	updatedAdmin, err := a.AdminRepo.UpdateAdmin(ctx, tx, admin)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("failed to update admin: %v", err)
	}

	return util.ConvertAdminToResponseDTO(updatedAdmin), nil
}