package service

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
)

type IuranService interface {
	AddIuran(iuran dto.IuranRequest) (string, error)
	UpdateIuran(id int64, iuran dto.IuranRequest) error
	GetAllIuran() ([]dto.IuranResponse, error)
	GetIuranById(id int64) (dto.IuranResponse, error)
	DeleteIuran(id int64) error
}

type iuranService struct {
	iuranRepo repository.IuranRepository
}

func NewIuranService(repo repository.IuranRepository) IuranService {
	return &iuranService{iuranRepo: repo}
}

func (s *iuranService) AddIuran(iuran dto.IuranRequest) (string, error) {
	// Generate UUID baru
	newID := uuid.New().String()
	iuran.ID = newID // pastikan dto.IuranRequest punya field ID bertipe string
	err := s.iuranRepo.Create(iuran)
	if err != nil {
		return "", err
	}
	return newID, nil
}

func (s *iuranService) UpdateIuran(id int64, iuran dto.IuranRequest) error {
	return s.iuranRepo.Update(id, iuran)
}

func (s *iuranService) GetAllIuran() ([]dto.IuranResponse, error) {
	return s.iuranRepo.GetAll()
}

func (s *iuranService) GetIuranById(id int64) (dto.IuranResponse, error) {
	return s.iuranRepo.GetByID(id)
}

func (s *iuranService) DeleteIuran(id int64) error {
	return s.iuranRepo.Delete(id)
}
