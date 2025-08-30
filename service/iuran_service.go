package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

// Konstanta untuk business logic iuran
const (
	IURAN_NOMINAL_PENUH = 7000 // Nominal penuh untuk status lunas
	STATUS_LUNAS        = "lunas"
	STATUS_BELUM_LUNAS  = "belum"
)

type IuranService interface {
	CreateMember(ctx context.Context, memberReq dto.MemberRequest) (dto.MemberResponse, int, error)
	GetAllMember(ctx context.Context) ([]dto.MemberResponse, int, error)
	GetMemberById(ctx context.Context, id string) (dto.MemberResponse, int, error)
	UpdateIuran(ctx context.Context, pembayaranReq dto.IuranRequest, id_member string) (dto.IuranResponse, int, error)
	DeleteMember(ctx context.Context, id_member string) (int, error)
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

	tanggalBayar, err := time.Parse("2006-01-02", pembayaranReq.TanggalBayar)
	if err != nil {
		return dto.IuranResponse{}, http.StatusBadRequest, fmt.Errorf("invalid tanggal_bayar format: %v", err)
	}

	// Validasi status iuran
	if pembayaranReq.Status != STATUS_LUNAS && pembayaranReq.Status != STATUS_BELUM_LUNAS {
		return dto.IuranResponse{}, http.StatusBadRequest, fmt.Errorf("status must be '%s' or '%s'", STATUS_LUNAS, STATUS_BELUM_LUNAS)
	}

	// Validasi jumlah bayar untuk status belum lunas
	if pembayaranReq.Status == STATUS_BELUM_LUNAS && pembayaranReq.JumlahBayar <= 0 {
		return dto.IuranResponse{}, http.StatusBadRequest, fmt.Errorf("jumlah_bayar harus lebih dari 0 untuk status '%s'", STATUS_BELUM_LUNAS)
	}

	// Check if a pembayaran already exists for this member, period, and week
	getPembayaran, err := i.IuranRepo.GetIuranByPeriod(ctx, tx, pembayaranReq.Periode, pembayaranReq.MingguKe, id_member)
	if err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get iuran by period: %v", err)
	}

	// Fetch member (assuming it's needed for validation or additional data)
	getMember, err := i.IuranRepo.GetMemberById(ctx, tx, id_member)
	if err != nil {
		return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get member by ID: %v", err)
	}
	// Use getMember if needed, e.g., to check if the member exists
	if getMember.IdMember == "" {
		return dto.IuranResponse{}, http.StatusBadRequest, fmt.Errorf("member not found")
	}

	var iuran model.Iuran
	var pembayaran model.PembayaranIuran // Declare outside to use in return
	if len(getPembayaran) == 0 {
		// No existing pembayaran, create a new iuran and pembayaran
		getIuran, err := i.IuranRepo.GetIuranByPeriodOnly(ctx, tx, pembayaranReq.Periode, pembayaranReq.MingguKe)
		if err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to check existing iuran: %v", err)
		}

		if getIuran == (model.Iuran{}) {
			newIuran := model.Iuran{
				IdIuran:  sql.NullString{String: uuid.New().String(), Valid: true},
				Periode:  sql.NullString{String: pembayaranReq.Periode, Valid: true},
				MingguKe: sql.NullInt64{Int64: int64(pembayaranReq.MingguKe), Valid: true},
			}
			iuran, err = i.IuranRepo.AddIuran(ctx, tx, newIuran)
			if err != nil {
				return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to add iuran: %v", err)
			}
		} else {
			iuran = getIuran
		}

		pembayaran = model.PembayaranIuran{
			IdPembayaran: sql.NullString{String: uuid.New().String(), Valid: true},
			IdMember:     sql.NullString{String: id_member, Valid: true},
			IdPemasukan:  sql.NullString{String: uuid.New().String(), Valid: true},
			Status:       sql.NullString{String: pembayaranReq.Status, Valid: true},
			TanggalBayar: sql.NullTime{Time: tanggalBayar, Valid: true},
			Iuran:        iuran,
		}
		// PERBAIKAN: Logic business untuk nilai iuran
		if pembayaranReq.Status == STATUS_LUNAS {
			// Jika status lunas, otomatis isi dengan nominal penuh iuran
			pembayaran.JumlahBayar = sql.NullInt64{Int64: IURAN_NOMINAL_PENUH, Valid: true}
		} else {
			// Jika status belum lunas, gunakan nilai yang diinput user
			pembayaran.JumlahBayar = sql.NullInt64{Int64: pembayaranReq.JumlahBayar, Valid: true}
		}

		_, err = i.IuranRepo.AddPembayaran(ctx, tx, pembayaran, getMember)
		if err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to add pembayaran: %v", err)
		}
		if err := tx.Commit(); err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
		}
		return util.ConvertIuranToResponseDTO(pembayaran), http.StatusCreated, nil
	} else {
		getPemb, err := i.IuranRepo.GetPembayaranById(ctx, tx, pembayaran, getMember.IdMember)
		if err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, err
		}

		// Update existing pembayaran
		existing := getPembayaran[0]
		pembayaran = model.PembayaranIuran{
			IdPembayaran: existing.IdPembayaran,
			IdMember:     sql.NullString{String: id_member, Valid: true},
			IdPemasukan:  getPemb.IdPemasukan,
			Status:       sql.NullString{String: pembayaranReq.Status, Valid: true},
			TanggalBayar: sql.NullTime{Time: tanggalBayar, Valid: true},
			Iuran: model.Iuran{
				IdIuran:  existing.Iuran.IdIuran,
				Periode:  existing.Iuran.Periode,
				MingguKe: sql.NullInt64{Int64: int64(pembayaranReq.MingguKe), Valid: true},
			},
		}

		// PERBAIKAN: Logic business untuk nilai iuran
		if pembayaranReq.Status == STATUS_LUNAS {
			// Jika status lunas, otomatis isi dengan nominal penuh iuran
			pembayaran.JumlahBayar = sql.NullInt64{Int64: IURAN_NOMINAL_PENUH, Valid: true}
		} else {
			// Jika status belum lunas, gunakan nilai yang diinput user
			pembayaran.JumlahBayar = sql.NullInt64{Int64: pembayaranReq.JumlahBayar, Valid: true}
		}

		updatedPembayaran, err := i.IuranRepo.UpdateStatusIuran(ctx, tx, pembayaran, getMember)
		if err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to update iuran status: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return dto.IuranResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
		}
		return util.ConvertIuranToResponseDTO(updatedPembayaran), http.StatusOK, nil
	}
}

// DeleteMember implements IuranService.
func (i *iuranService) DeleteMember(ctx context.Context, id_member string) (int, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	member, err := i.IuranRepo.GetMemberById(ctx, tx, id_member)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get member by ID: %v", err)
	}

	if member.IdMember == "" {
		return http.StatusBadRequest, fmt.Errorf("member not found")
	}

	err = i.IuranRepo.DeleteMember(ctx, tx, id_member)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to delete member: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return http.StatusNoContent, nil
}
