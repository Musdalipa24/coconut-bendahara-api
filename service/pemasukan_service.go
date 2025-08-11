package service

import (
	"context"
	"net/http"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

type PemasukanService interface {
	AddPemasukan(ctx context.Context, r *http.Request, pemasukanRequest dto.PemasukanRequest) (dto.PemasukanResponse, error)
	UpdatePemasukan(ctx context.Context, r *http.Request, pemasukanRequest dto.PemasukanRequest, no string) (dto.PemasukanResponse, error)
	GetById(ctx context.Context, id string) (dto.PemasukanResponse, error)
	DeletePemasukan(ctx context.Context, id string) (dto.PemasukanResponse, error)
	GetPemasukan(ctx context.Context, page int, pageSize int) (dto.PemasukanPaginationResponse, error)
	GetPemasukanByDateRange(ctx context.Context, startDate, endDate string, page int, pageSize int) (dto.PemasukanPaginationResponse, error)
}
