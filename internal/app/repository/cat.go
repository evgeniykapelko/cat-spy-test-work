package repository

import (
	"gorm.io/gorm"
	"log"
	"spy_cat/internal/app/model"
)

type CatRepository interface {
	Create(cat *model.Cat) error
	FindAll() ([]model.Cat, error)
	FindByID(id uint) (*model.Cat, error)
	Update(cat *model.Cat) error
	Delete(cat *model.Cat) error
}

type catRepository struct {
	DB *gorm.DB
}

func NewCatRepository(db *gorm.DB) CatRepository {
	if db == nil {
		log.Fatal("DB instance is nil")
	}

	return &catRepository{DB: db}
}

func (r *catRepository) Create(cat *model.Cat) error {
	return r.DB.Create(cat).Error
}

func (r *catRepository) FindAll() ([]model.Cat, error) {
	var cats []model.Cat
	err := r.DB.Find(&cats).Error
	return cats, err
}

func (r *catRepository) FindByID(id uint) (*model.Cat, error) {
	var cat model.Cat
	err := r.DB.First(&cat, id).Error
	return &cat, err
}

func (r *catRepository) Update(cat *model.Cat) error {
	return r.DB.Save(cat).Error
}

func (r *catRepository) Delete(cat *model.Cat) error {
	return r.DB.Delete(cat).Error
}
