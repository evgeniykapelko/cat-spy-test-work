package service

import (
	"spy_cat/internal/app/model"
	"spy_cat/internal/app/repository"
)

type TargetService struct {
	repo repository.TargetRepository
}

func NewTargetService(repo repository.TargetRepository) *TargetService {
	return &TargetService{repo: repo}
}

func (s *TargetService) Create(target *model.Target) error {
	return s.repo.Create(target)
}

func (s *TargetService) FindAll() ([]model.Target, error) {
	return s.repo.FindAll()
}

func (s *TargetService) FindByID(id uint) (*model.Target, error) {
	return s.repo.FindByID(id)
}

func (s *TargetService) Update(target *model.Target) error {
	return s.repo.Update(target)
}

func (s *TargetService) Delete(target *model.Target) error {
	return s.repo.Delete(target)
}

func (s *TargetService) CompleteTarget(targetID uint) error {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return err
	}
	target.Complete = true
	return s.repo.Update(target)
}
