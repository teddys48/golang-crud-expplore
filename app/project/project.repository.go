package project

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Repository interface {
	Create(db *gorm.DB, user *Project) error
	Update(db *gorm.DB, user *Project, id int64) error
	CheckByID(db *gorm.DB, user *Project, id int64) error
	Delete(db *gorm.DB, user *Project, id int64) error
	All(db *gorm.DB, user *[]Project) error
}

type repository struct {
}

func Newrepository(config *viper.Viper) Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, data *Project) error {
	return db.Table("project").Create(&data).Error
}

func (r *repository) Update(db *gorm.DB, data *Project, id int64) error {
	return db.Table("project").
		Where("id", id).
		Updates(&data).Error
}

func (r *repository) CheckByID(db *gorm.DB, data *Project, id int64) error {
	return db.Table("project").
		Where("id", id).
		Where("status", "ACTIVE").
		First(&data).
		Scan(data).Error
}

func (r *repository) Delete(db *gorm.DB, data *Project, id int64) error {
	return db.Table("project").Delete(&data).Error
}

func (r *repository) All(db *gorm.DB, data *[]Project) error {
	return db.Table("project").Find(&data).Error
}
