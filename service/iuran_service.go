package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type IuranService interface {
	CreateMember(ctx context.Context, memberReq dto.MemberRequest) (dto.MemberResponse, int, error)
	GetAllMember(ctx context.Context) ([]dto.MemberResponse, int, error)
	GetMemberById(ctx context.Context, id string) (dto.MemberResponse, int, error)
	UpdateIuran(ctx context.Context, pembayaranReq dto.IuranRequest, id_member string) (dto.IuranResponse, int, error)
}

type iuranService struct {
	IuranRepo repository.IuranRepository
	DB        *sql.DB
}

func NewIuranService(iuranRepo repository.IuranRepository, db *sql.DB) IuranService {
	return &iuranService{
		IuranRepo: iuranRepo,
		DB:        db,
	}
}

// CreateMember implements IuranService.
func (i *iuranService) CreateMember(ctx context.Context, memberReq dto.MemberRequest) (dto.MemberResponse, int, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	member := model.Member{
		IdMember: uuid.New().String(),
		NRA:      memberReq.NRA,
		Nama:     memberReq.Nama,
		Status:   memberReq.Status,
	}

	addedMember, err := i.IuranRepo.AddMember(ctx, tx, member)
	if err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to add member: %v", err)
	}

	getMember, err := i.IuranRepo.GetMemberById(ctx, tx, addedMember.IdMember)
	if err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get member: %v", err)
	}

	addedMember.CreatedAt = getMember.CreatedAt
	addedMember.UpdatedAt = getMember.UpdatedAt

	if err := tx.Commit(); err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return util.ConvertMemberToResponseDTO(addedMember), http.StatusCreated, nil
}

// GetAllMember implements MemberService.
func (i *iuranService) GetAllMember(ctx context.Context) ([]dto.MemberResponse, int, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return []dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	members, err := i.IuranRepo.GetAllMembers(ctx, tx)
	if err != nil {
		return []dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get all members: %v", err)
	}

	for idx := range members {
		iuran, err := i.IuranRepo.GetIuranByMemberID(ctx, tx, members[idx].IdMember)
		if err != nil {
			return []dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get all iuran: %v", err)
		}

		members[idx].PembayaranIurans = append(members[idx].PembayaranIurans, iuran...)
	}

	if err := tx.Commit(); err != nil {
		return []dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return util.ConvertMemberToListResponseDTO(members), http.StatusOK, nil
}

// GetMemberById implements MemberService.
func (i *iuranService) GetMemberById(ctx context.Context, id string) (dto.MemberResponse, int, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	getMember, err := i.IuranRepo.GetMemberById(ctx, tx, id)
	if err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get member: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return dto.MemberResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return util.ConvertMemberToResponseDTO(getMember), http.StatusOK, nil
}

// UpdateIuran implements IuranService.
func (i *iuranService) UpdateIuran(ctx context.Context, pembayaranReq dto.IuranRequest, id_member string) (dto.IuranResponse, int, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Ambil pembayaran di periode+minggu_ke ini
	getPembayaran, err := i.IuranRepo.GetIuranByPeriod(ctx, tx, pembayaranReq.Periode, pembayaranReq.MingguKe)
	if err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get iuran by period: %v", err)
	}

	getMember, err := i.IuranRepo.GetMemberById(ctx, tx, id_member)
	if err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get member by ID: %v", err)
	}

	// =========================
	// CASE 1: Belum ada pembayaran
	// =========================
	if len(getPembayaran) == 0 || getMember.IdMember == id_member {
		// Cek apakah iuran di periode ini sudah ada
		getIuran, err := i.IuranRepo.GetIuranByPeriodOnly(ctx, tx, pembayaranReq.Periode, pembayaranReq.MingguKe)
		if err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to check existing iuran: %v", err)
		}

		var iuran model.Iuran
		if getIuran == (model.Iuran{}) {
			// Buat iuran baru
			newIuran := model.Iuran{
				IdIuran:  sql.NullString{String: uuid.New().String(), Valid: true},
				Periode:  sql.NullString{String: pembayaranReq.Periode, Valid: true},
				MingguKe: sql.NullInt64{Int64: int64(pembayaranReq.MingguKe), Valid: true},
			}

			iuran, err = i.IuranRepo.AddIuran(ctx, tx, newIuran)
			if err != nil {
				return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to add iuran: %v", err)
			}

			fmt.Println("Iuran baru dibuat:", iuran)
		} else {
			// Gunakan iuran yang sudah ada
			iuran = getIuran
			fmt.Println("Menggunakan iuran yang sudah ada:", iuran)
		}

		// Buat pembayaran baru
		pembayaran := model.PembayaranIuran{
			IdPembayaran: sql.NullString{String: uuid.New().String(), Valid: true},
			IdMember:     sql.NullString{String: id_member, Valid: true},
			Status:       sql.NullString{String: pembayaranReq.Status, Valid: true},
			TanggalBayar: sql.NullString{String: pembayaranReq.TanggalBayar, Valid: true},
			Iuran:        iuran,
		}

		createdPembayaran, err := i.IuranRepo.AddPembayaran(ctx, tx, pembayaran)
		if err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to add pembayaran: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
		}

		fmt.Println("Pembayaran baru dibuat:", createdPembayaran)
		return util.ConvertIuranToResponseDTO(createdPembayaran), http.StatusCreated, nil
	}

	// =========================
	// CASE 2: Sudah ada pembayaran → UPDATE
	// =========================
	existing := getPembayaran[0] // aman karena len > 0

	pembayaran := model.PembayaranIuran{
		IdPembayaran: existing.IdPembayaran,
		IdMember:     sql.NullString{String: id_member, Valid: true},
		Status:       sql.NullString{String: pembayaranReq.Status, Valid: true},
		TanggalBayar: sql.NullString{String: pembayaranReq.TanggalBayar, Valid: true},
		Iuran: model.Iuran{
			IdIuran:  existing.Iuran.IdIuran,
			Periode:  existing.Iuran.Periode,
			MingguKe: sql.NullInt64{Int64: int64(pembayaranReq.MingguKe), Valid: true},
		},
	}

	updatedPembayaran, err := i.IuranRepo.UpdateStatusIuran(ctx, tx, pembayaran)
	if err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to update iuran status: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Println("Periode", pembayaranReq.Periode, "minggu ke", pembayaranReq.MingguKe, "telah diperbarui", "id pembayaran:", updatedPembayaran.IdPembayaran.String, "status:", updatedPembayaran.Status.String, "id member:", updatedPembayaran.IdMember.String)

	return util.ConvertIuranToResponseDTO(updatedPembayaran), http.StatusOK, nil
}