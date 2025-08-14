package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertMemberToResponseDTO(member model.Member) dto.MemberResponse {
	var iuranResponses []dto.IuranResponse
	for _, iuran := range member.PembayaranIurans {
		iuranResponses = append(iuranResponses, ConvertIuranToResponseDTO(iuran))
	}

	return dto.MemberResponse{
		IdMember:  member.IdMember,
		NRA:       member.NRA,
		Nama:      member.Nama,
		Status:    member.Status,
		Iuran:     iuranResponses,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func ConvertMemberToListResponseDTO(member []model.Member) []dto.MemberResponse {
	var memberResponse []dto.MemberResponse
	for _, members := range member {
		memberResponse = append(memberResponse, ConvertMemberToResponseDTO(members))
	}

	return memberResponse
}
