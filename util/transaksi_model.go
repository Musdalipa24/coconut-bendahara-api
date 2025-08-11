package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertTransactionToResponseDTO(transaction model.Transaction) dto.TransactionResponse {
	// Format tanggal sesuai dengan format yang diinginkan (misalnya: "02-01-2006 15:04")
	formattedTanggal := transaction.Tanggal.Format("02-01-2006 15:04")

	return dto.TransactionResponse{
		Id:             transaction.Id,
		Tanggal:        formattedTanggal,
		Keterangan:     transaction.Keterangan,
		JenisTransaksi: transaction.JenisTransaksi,
		Nominal:        transaction.Nominal,
	}
}

func ConvertTransactionToListResponseDTO(transaction []model.Transaction) []dto.TransactionResponse {
	var transactionResponse []dto.TransactionResponse
	for _, transactions := range transaction {
		transactionResponse = append(transactionResponse, ConvertTransactionToResponseDTO(transactions))
	}

	return transactionResponse
}
