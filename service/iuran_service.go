package service

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
)

type IuranService interface {
	AddIuran(iuran dto.IuranRequest) (string, error)
	UpdateIuran(id string, iuran dto.IuranRequest) error
	GetAllIuran() ([]dto.IuranResponse, error)
	GetIuranById(id string) (dto.IuranResponse, error)
	DeleteIuran(id int64) error
}

type iuranService struct {
	iuranRepo repository.IuranRepository
}

func NewIuranService(repo repository.IuranRepository) IuranService {
	return &iuranService{iuranRepo: repo}
}

func (s *iuranService) AddIuran(iuran dto.IuranRequest) (string, error) {
	// Karena ID bertipe int64 dan biasanya auto increment di DB,
	// cukup panggil Create(iuran), lalu ambil ID dari DB jika perlu.

	err := s.iuranRepo.Create(iuran)
	if err != nil {
		return "", err
	}
	// Jika ingin mengembalikan ID, ambil dari DB setelah insert
	return "", nil
}

func (s *iuranService) UpdateIuran(id string, iuran dto.IuranRequest) error {
	return s.iuranRepo.Update(id, iuran)
}

func (s *iuranService) GetAllIuran() ([]dto.IuranResponse, error) {
	return s.iuranRepo.GetAll()
}

func (s *iuranService) GetIuranById(id string) (dto.IuranResponse, error) {
	return s.iuranRepo.GetByID(id)
}

func (s *iuranService) DeleteIuran(id int64) error {
	return s.iuranRepo.Delete(id)
}
