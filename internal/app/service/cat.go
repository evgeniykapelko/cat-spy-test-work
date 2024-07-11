package service

import (
	"spy_cat/internal/app/model"
	"spy_cat/internal/app/repository"
)

type CatService struct {
	repo           repository.CatRepository
	breedValidator *BreedValidator
}

func NewCatService(repo repository.CatRepository) *CatService {
	return &CatService{
		repo:           repo,
		breedValidator: NewBreedValidator(),
	}
}

func (s *CatService) Create(cat *model.Cat) error {
	return s.repo.Create(cat)
}

func (s *CatService) FindAll() ([]model.Cat, error) {
	return s.repo.FindAll()
}

func (s *CatService) FindByID(id uint) (*model.Cat, error) {
	return s.repo.FindByID(id)
}

func (s *CatService) Update(cat *model.Cat) error {
	return s.repo.Update(cat)
}

func (s *CatService) Delete(cat *model.Cat) error {
	return s.repo.Delete(cat)
}
