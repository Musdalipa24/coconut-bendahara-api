package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type TransactionService interface {
	GetAllTransaction(ctx context.Context) ([]dto.TransactionResponse, int, error)
	GetLastTransaction(ctx context.Context) ([]dto.TransactionResponse, int, error)
}

type transactionServiceImpl struct {
	TransactionRepo repository.TransaksiRepo
	DB              *sql.DB
}

func NewTransactionService(transactionRepo repository.TransaksiRepo, db *sql.DB) TransactionService {
	return &transactionServiceImpl{
		TransactionRepo: transactionRepo,
		DB:              db,
	}
}

// GetAllTransaction implements TransactionService.
func (t *transactionServiceImpl) GetAllTransaction(ctx context.Context) ([]dto.TransactionResponse, int, error) {
	tx, err := t.DB.Begin()
	if err != nil {
		return []dto.TransactionResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	transaction, err := t.TransactionRepo.GetAllTransaction(ctx, tx)
	if err != nil {
		return []dto.TransactionResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get all transaction")
	}

	return util.ConvertTransactionToListResponseDTO(transaction), http.StatusOK, nil
}

// GetLastTransaction implements TransactionService.
func (t *transactionServiceImpl) GetLastTransaction(ctx context.Context) ([]dto.TransactionResponse, int, error) {
	tx, err := t.DB.Begin()
	if err != nil {
		return []dto.TransactionResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction")
	}
	defer util.CommitOrRollBack(tx)

	transaction, err := t.TransactionRepo.GetLastTransaction(ctx, tx)
	if err != nil {
		return []dto.TransactionResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get last transaction")
	}

	return util.ConvertTransactionToListResponseDTO(transaction), http.StatusOK, nil
}
