package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertAdminToResponseDTO(admin model.Admin) dto.AdminResponse {
	return dto.AdminResponse{
		Id:       admin.Id,
		Nik:      admin.Nik,
		Username: admin.Username,
		Password: admin.Password,
		Role:     admin.Role,
	}
}

func ConvertAdminToListResponseDTO(admin []model.Admin) []dto.AdminResponse {
	var adminResponse []dto.AdminResponse
	for _, admins := range admin {
		adminResponse = append(adminResponse, ConvertAdminToResponseDTO(admins))
	}

	return adminResponse
}
