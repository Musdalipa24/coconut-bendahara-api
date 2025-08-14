package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type PemasukanService interface {
	AddPemasukan(ctx context.Context, r *http.Request, pemasukanRequest dto.PemasukanRequest) (dto.PemasukanResponse, error)
	UpdatePemasukan(ctx context.Context, r *http.Request, pemasukanRequest dto.PemasukanRequest, no string) (dto.PemasukanResponse, error)
	GetById(ctx context.Context, id string) (dto.PemasukanResponse, error)
	DeletePemasukan(ctx context.Context, id string) (dto.PemasukanResponse, error)
	GetPemasukan(ctx context.Context, page int, pageSize int) (dto.PemasukanPaginationResponse, error)
	GetPemasukanByDateRange(ctx context.Context, startDate, endDate string, page int, pageSize int) (dto.PemasukanPaginationResponse, error)
}

type pemasukanServiceImpl struct {
	PemasukanRepo repository.PemasukanRepo
	DB            *sql.DB
}

func NewPemasukanService(pemasukanRepo repository.PemasukanRepo, db *sql.DB) PemasukanService {
	return &pemasukanServiceImpl{
		PemasukanRepo: pemasukanRepo,
		DB:            db,
	}
}

// AddPemasukan implements PemasukanService.
func (s *pemasukanServiceImpl) AddPemasukan(ctx context.Context, r *http.Request, pemasukanRequest dto.PemasukanRequest) (dto.PemasukanResponse, error) {
	// Parse multipart form dengan batas ukuran 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to parse form: %v", err)
	}

	// Ambil nilai dari form
	tanggalStr := r.FormValue("tanggal") // Format: "2006-01-02 15:04"
	kategori := r.FormValue("kategori")
	keterangan := r.FormValue("keterangan")
	nominalStr := r.FormValue("nominal")

	// Parse tanggal dari string ke time.Time
	tanggal, err := time.Parse("2006-01-02 15:04", tanggalStr)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("invalid date format, expected 'YYYY-MM-DD HH:MM': %v", err)
	}

	// Konversi nominal dari string ke uint64
	nominal, err := strconv.ParseUint(nominalStr, 10, 64)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("nominal must be a valid number: %v", err)
	}

	// Buat DTO request
	pemasukanRequest = dto.PemasukanRequest{
		Tanggal:    tanggalStr,
		Kategori:   kategori,
		Keterangan: keterangan,
		Nominal:    nominal,
	}

	// Ambil file dari form (opsional)
	var fileName string
	file, handler, err := r.FormFile("nota")
	if err == nil && file != nil {
		defer file.Close()

		// Format tanggal untuk nama file
		formattedDateTime := tanggal.Format("2006-01-02-15-04") // Format: YYYY-MM-DD-HH-MM

		// Buat nama file dengan format: tanggal-waktu-uuid
		fileName = fmt.Sprintf("%s-%s%s", formattedDateTime, uuid.New().String(), filepath.Ext(handler.Filename))

		// Buat direktori upload jika belum ada
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// Simpan file ke direktori upload
		filePath := filepath.Join(uploadDir, fileName)
		out, err := os.Create(filePath)
		if err != nil {
			return dto.PemasukanResponse{}, fmt.Errorf("failed to create file: %v", err)
		}
		defer out.Close()

		// Salin file yang diunggah ke file W yang baru dibuat
		_, err = io.Copy(out, file)
		if err != nil {
			return dto.PemasukanResponse{}, fmt.Errorf("failed to copy file: %v", err)
		}

		// Simpan nama file ke dalam request
		pemasukanRequest.Nota = fileName
	}

	// Validasi input (Nota sekarang opsional)
	if pemasukanRequest.Tanggal == "" || pemasukanRequest.Kategori == "" || pemasukanRequest.Nominal == 0 {
		return dto.PemasukanResponse{}, fmt.Errorf("date, category, or nominal can't be empty")
	}

	// Mulai transaksi database
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Buat objek Pemasukan
	pemasukan := model.Pemasukan{
		Id:         uuid.New().String(),
		Tanggal:    tanggal,
		Kategori:   pemasukanRequest.Kategori,
		Keterangan: pemasukanRequest.Keterangan,
		Nominal:    pemasukanRequest.Nominal,
		Nota:       pemasukanRequest.Nota,
	}

	// Tambahkan pemasukan ke database
	addPemasukan, err := s.PemasukanRepo.AddPemasukan(ctx, tx, pemasukan)
	if err != nil {
		// Hapus file yang sudah diunggah jika transaksi gagal dan file ada
		if fileName != "" {
			os.Remove(filepath.Join("./uploads", fileName))
		}
		return dto.PemasukanResponse{}, fmt.Errorf("failed to add income: %v", err)
	}

	// Kembalikan respons
	return util.ConvertPemasukanToResponseDTO(addPemasukan), nil
}

// UpdatePemasukan implements PemasukanService.
func (s *pemasukanServiceImpl) UpdatePemasukan(ctx context.Context, r *http.Request, pemasukanRequest dto.PemasukanRequest, id string) (dto.PemasukanResponse, error) {
	// Parse multipart form dengan batas ukuran 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to parse form: %v", err)
	}

	// Ambil nilai dari form
	tanggalStr := r.FormValue("tanggal") // Format: "2006-01-02 15:04"
	kategori := r.FormValue("kategori")
	keterangan := r.FormValue("keterangan")
	nominalStr := r.FormValue("nominal")

	// Parse tanggal dari string ke time.Time
	tanggal, err := time.Parse("2006-01-02 15:04", tanggalStr)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("invalid date format, expected 'YYYY-MM-DD HH:MM': %v", err)
	}

	// Konversi nominal dari string ke uint64
	nominal, err := strconv.ParseUint(nominalStr, 10, 64)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("nominal must be a valid number: %v", err)
	}

	// Buat DTO request
	pemasukanRequest = dto.PemasukanRequest{
		Tanggal:    tanggalStr,
		Kategori:   kategori,
		Keterangan: keterangan,
		Nominal:    nominal,
	}

	// Ambil file dari form (opsional)
	var fileName string
	file, handler, err := r.FormFile("nota")
	if err == nil && file != nil {
		defer file.Close()

		// Format tanggal untuk nama file
		formattedDateTime := tanggal.Format("2006-01-02-15-04") // Format: YYYY-MM-DD-HH-MM

		// Buat nama file dengan format: tanggal-waktu-uuid
		fileName = fmt.Sprintf("%s-%s%s", formattedDateTime, uuid.New().String(), filepath.Ext(handler.Filename))

		// Buat direktori upload jika belum ada
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// Simpan file ke direktori upload
		filePath := filepath.Join(uploadDir, fileName)
		out, err := os.Create(filePath)
		if err != nil {
			return dto.PemasukanResponse{}, fmt.Errorf("failed to create file: %v", err)
		}
		defer out.Close()

		// Salin file yang diunggah ke file W yang baru dibuat
		_, err = io.Copy(out, file)
		if err != nil {
			return dto.PemasukanResponse{}, fmt.Errorf("failed to copy file: %v", err)
		}

		// Simpan nama file ke dalam request
		pemasukanRequest.Nota = fileName
	}

	// Validasi input (Nota sekarang opsional)
	if pemasukanRequest.Tanggal == "" || pemasukanRequest.Kategori == "" || pemasukanRequest.Nominal == 0 {
		return dto.PemasukanResponse{}, fmt.Errorf("date, category, or nominal can't be empty")
	}

	// Mulai transaksi database
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Buat objek Pemasukan
	pemasukan := model.Pemasukan{
		Id:         uuid.New().String(),
		Tanggal:    tanggal,
		Kategori:   pemasukanRequest.Kategori,
		Keterangan: pemasukanRequest.Keterangan,
		Nominal:    pemasukanRequest.Nominal,
		Nota:       pemasukanRequest.Nota,
	}

	// Simpan perubahan ke database
	updatePemasukan, err := s.PemasukanRepo.UpdatePemasukan(ctx, tx, pemasukan, id)
	if err != nil {
		panic(err)
		// return dto.PemasukanResponse{}, fmt.Errorf("failed to update income")
	}

	return util.ConvertPemasukanToResponseDTO(updatePemasukan), nil
}

// GetPemasukan implements PemasukanService.
func (s *pemasukanServiceImpl) GetPemasukan(ctx context.Context, page int, pageSize int) (dto.PemasukanPaginationResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PemasukanPaginationResponse{}, fmt.Errorf("failed to start transaction")
	}
	defer tx.Commit()

	// Validasi parameter pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15 // Default page size
	}

	pemasukan, total, err := s.PemasukanRepo.GetPemasukan(ctx, tx, page, pageSize)
	if err != nil {
		return dto.PemasukanPaginationResponse{}, err
	}

	// Hitung total halaman
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	response := dto.PemasukanPaginationResponse{
		Items:      util.ConvertPemasukanToListResponseDTO(pemasukan),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return response, nil
}

// DeletePemasukan implements PemasukanService.
func (s *pemasukanServiceImpl) DeletePemasukan(ctx context.Context, id string) (dto.PemasukanResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to start transaction")
	}
	defer tx.Commit()

	pemasukan, err := s.PemasukanRepo.FindById(ctx, tx, id)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("income not found")
	}

	pemasukan, err = s.PemasukanRepo.DeletePemasukan(ctx, tx, pemasukan)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to delete income")
	}

	return util.ConvertPemasukanToResponseDTO(pemasukan), nil
}

// GetById implements PemasukanService.
func (s *pemasukanServiceImpl) GetById(ctx context.Context, id string) (dto.PemasukanResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("failed to start transaction")
	}
	defer tx.Commit()

	pemasukan, err := s.PemasukanRepo.FindById(ctx, tx, id)
	if err != nil {
		return dto.PemasukanResponse{}, fmt.Errorf("income not found")
	}

	return util.ConvertPemasukanToResponseDTO(pemasukan), nil
}

// GetPemasukanByDateRange implements PemasukanService.
func (s *pemasukanServiceImpl) GetPemasukanByDateRange(ctx context.Context, startDate, endDate string, page int, pageSize int) (dto.PemasukanPaginationResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PemasukanPaginationResponse{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Validasi parameter pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15 // Default page size
	}

	// Validasi tanggal
	if startDate == "" || endDate == "" {
		return dto.PemasukanPaginationResponse{}, fmt.Errorf("start_date and end_date are required")
	}

	// Parse tanggal untuk memastikan format valid
	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		return dto.PemasukanPaginationResponse{}, fmt.Errorf("invalid start_date format, expected 'YYYY-MM-DD': %v", err)
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return dto.PemasukanPaginationResponse{}, fmt.Errorf("invalid end_date format, expected 'YYYY-MM-DD': %v", err)
	}

	pemasukan, total, err := s.PemasukanRepo.GetPemasukanByDateRange(ctx, tx, startDate, endDate, page, pageSize)
	if err != nil {
		return dto.PemasukanPaginationResponse{}, err
	}

	// Hitung total halaman
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	response := dto.PemasukanPaginationResponse{
		Items:      util.ConvertPemasukanToListResponseDTO(pemasukan),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return response, nil
}
