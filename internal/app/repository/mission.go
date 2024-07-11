package repository

import (
	"gorm.io/gorm"
	"log"
	"spy_cat/internal/app/model"
)

type MissionRepository interface {
	Create(mission *model.Mission) error
	FindAll() ([]model.Mission, error)
	FindByID(id uint) (*model.Mission, error)
	Update(mission *model.Mission) error
	Delete(mission *model.Mission) error
}

type missionRepository struct {
	DB *gorm.DB
}

func NewMissionRepository(db *gorm.DB) MissionRepository {
	if db == nil {
		log.Fatal("DB instance is nil")
	}

	return &missionRepository{DB: db}
}

func (r *missionRepository) Create(mission *model.Mission) error {
	return r.DB.Create(mission).Error
}

func (r *missionRepository) FindAll() ([]model.Mission, error) {
	var missions []model.Mission
	err := r.DB.Preload("Targets").Find(&missions).Error
	return missions, err
}

func (r *missionRepository) FindByID(id uint) (*model.Mission, error) {
	var mission model.Mission
	err := r.DB.Preload("Targets").First(&mission, id).Error
	return &mission, err
}

func (r *missionRepository) Update(mission *model.Mission) error {
	return r.DB.Save(mission).Error
}

func (r *missionRepository) Delete(mission *model.Mission) error {
	return r.DB.Delete(mission).Error
}
