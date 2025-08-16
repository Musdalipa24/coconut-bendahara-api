package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertIuranToResponseDTO(pembayaran model.PembayaranIuran) dto.IuranResponse {
	return dto.IuranResponse{
		IdIuran:      pembayaran.Iuran.IdIuran.String,
		Periode:      pembayaran.Iuran.Periode.String,
		MingguKe:     int(pembayaran.Iuran.MingguKe.Int64),
		Status:       pembayaran.Status.String,
		TanggalBayar: pembayaran.TanggalBayar.Time.String(),
	}
}
