package service

import (
	"context"
	"net/http"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

type PengeluaranService interface {
	AddPengeluaran(ctx context.Context, r *http.Request, pengeluaranRequest dto.PengeluaranRequest) (dto.PengeluaranResponse, error)
	UpdatePengeluaran(ctx context.Context, r *http.Request, pengeluaranRequest dto.PengeluaranRequest, no string) (dto.PengeluaranResponse, error)
	GetById(ctx context.Context, id string) (dto.PengeluaranResponse, error)
	DeletePengeluaran(ctx context.Context, id string) (dto.PengeluaranResponse, error)
	GetPengeluaran(ctx context.Context, page int, pageSize int) (dto.PengeluaranPaginationResponse, error)
	GetPengeluaranByDateRange(ctx context.Context, startDate, endDate string, page int, pageSize int) (dto.PengeluaranPaginationResponse, error)
}
