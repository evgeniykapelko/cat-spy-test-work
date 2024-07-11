package repository

import (
	"gorm.io/gorm"
	"log"
	"spy_cat/internal/app/model"
)

type TargetRepository interface {
	Create(target *model.Target) error
	FindAll() ([]model.Target, error)
	FindByID(id uint) (*model.Target, error)
	Update(target *model.Target) error
	Delete(target *model.Target) error
	UpdateNotes(id uint, notes string) error
}

type targetRepository struct {
	DB *gorm.DB
}

func NewTargetRepository(db *gorm.DB) TargetRepository {
	if db == nil {
		log.Fatal("DB instance is nil")
	}

	return &targetRepository{DB: db}
}

func (r *targetRepository) Create(target *model.Target) error {
	return r.DB.Create(target).Error
}

func (r *targetRepository) FindAll() ([]model.Target, error) {
	var targets []model.Target
	err := r.DB.Find(&targets).Error
	return targets, err
}

func (r *targetRepository) FindByID(id uint) (*model.Target, error) {
	var target model.Target
	err := r.DB.First(&target, id).Error
	return &target, err
}

func (r *targetRepository) Update(target *model.Target) error {
	return r.DB.Save(target).Error
}

func (r *targetRepository) Delete(target *model.Target) error {
	return r.DB.Delete(target).Error
}

func (r *targetRepository) UpdateNotes(id uint, notes string) error {
	return r.DB.Model(&model.Target{}).Where("id = ?", id).Update("notes", notes).Error
}
