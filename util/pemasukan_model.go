package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertPemasukanToResponseDTO(pemasukan model.Pemasukan) dto.PemasukanResponse {
	// Format tanggal sesuai dengan format yang diinginkan (misalnya: "02-01-2006 15:04")
	formattedTanggal := pemasukan.Tanggal.Format("02-01-2006 15:04")

	return dto.PemasukanResponse{
		Id:         pemasukan.Id,
		Tanggal:    formattedTanggal,
		Kategori:   pemasukan.Kategori,
		Keterangan: pemasukan.Keterangan,
		Nominal:    pemasukan.Nominal,
		Nota:       pemasukan.Nota,
	}
}

func ConvertPemasukanToListResponseDTO(pemasukan []model.Pemasukan) []dto.PemasukanResponse {
	var pemasukanResponse []dto.PemasukanResponse
	for _, pemasukans := range pemasukan {
		pemasukanResponse = append(pemasukanResponse, ConvertPemasukanToResponseDTO(pemasukans))
	}

	return pemasukanResponse
}
