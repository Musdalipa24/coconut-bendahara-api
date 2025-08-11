package service

import (
	"context"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
)

type TransactionService interface {
	GetAllTransaction(ctx context.Context) ([]dto.TransactionResponse, error)
	GetLastTransaction(ctx context.Context) ([]dto.TransactionResponse, error)
}
