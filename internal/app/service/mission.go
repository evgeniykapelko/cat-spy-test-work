package service

import (
	"errors"
	"spy_cat/internal/app/model"
	"spy_cat/internal/app/repository"
)

type MissionService struct {
	repo       repository.MissionRepository
	catRepo    repository.CatRepository
	targetRepo repository.TargetRepository
}

func NewMissionService(repo repository.MissionRepository, catRepo repository.CatRepository, targetRepo repository.TargetRepository) *MissionService {
	return &MissionService{
		repo:       repo,
		catRepo:    catRepo,
		targetRepo: targetRepo,
	}
}

func (s *MissionService) Create(mission *model.Mission) error {
	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		return errors.New("mission must have between 1 and 3 targets")
	}
	return s.repo.Create(mission)
}

func (s *MissionService) Delete(mission *model.Mission) error {
	if mission.CatID != 0 {
		return errors.New("cannot delete mission assigned to a cat")
	}
	return s.repo.Delete(mission)
}

func (s *MissionService) AddTarget(mission *model.Mission, target *model.Target) error {
	if mission.Complete {
		return errors.New("cannot add target to a completed mission")
	}
	if len(mission.Targets) >= 3 {
		return errors.New("cannot add more than 3 targets to a mission")
	}
	mission.Targets = append(mission.Targets, *target)
	return s.repo.Update(mission)
}

func (s *MissionService) FindAll() ([]model.Mission, error) {
	return s.repo.FindAll()
}

func (s *MissionService) FindByID(id uint) (*model.Mission, error) {
	return s.repo.FindByID(id)
}

func (s *MissionService) Update(mission *model.Mission) error {
	return s.repo.Update(mission)
}

func (s *MissionService) AssignCatToMission(missionID uint, catID uint) error {
	mission, err := s.repo.FindByID(missionID)
	if err != nil {
		return err
	}
	if mission.Complete {
		return errors.New("cannot assign cat to a completed mission")
	}
	mission.CatID = catID
	return s.repo.Update(mission)
}

func (s *MissionService) CompleteMission(missionID uint) error {
	mission, err := s.repo.FindByID(missionID)
	if err != nil {
		return err
	}
	mission.Complete = true
	return s.repo.Update(mission)
}

func (s *MissionService) UpdateTargetNotes(missionId, targetId uint, notes string) error {
	return s.targetRepo.UpdateNotes(targetId, notes)
}

func (s *MissionService) FindTargetByID(id uint) (*model.Target, error) {
	return s.targetRepo.FindByID(id)
}

func (s *MissionService) DeleteTarget(target *model.Target) error {
	return s.targetRepo.Delete(target)
}
