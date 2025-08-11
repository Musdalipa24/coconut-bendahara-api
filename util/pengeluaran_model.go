package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertPengeluaranToResponseDTO(pengeluaran model.Pengeluaran) dto.PengeluaranResponse {
	// Format tanggal sesuai dengan format yang diinginkan (misalnya: "02-01-2006 15:04")
	formattedTanggal := pengeluaran.Tanggal.Format("02-01-2006 15:04")

	return dto.PengeluaranResponse{
		Id:         pengeluaran.Id,
		Tanggal:    formattedTanggal,
		Nota:       pengeluaran.Nota,
		Nominal:    pengeluaran.Nominal,
		Keterangan: pengeluaran.Keterangan,
	}
}

func ConvertPengeluaranToListResponseDTO(pengeluaran []model.Pengeluaran) []dto.PengeluaranResponse {
	var pengeluaranResponse []dto.PengeluaranResponse
	for _, pengeluarans := range pengeluaran {
		pengeluaranResponse = append(pengeluaranResponse, ConvertPengeluaranToResponseDTO(pengeluarans))
	}

	return pengeluaranResponse
}
