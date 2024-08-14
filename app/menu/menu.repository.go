package menu

import (
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(db *gorm.DB, user *entity.Menu) error
	Update(db *gorm.DB, user *entity.Menu, id int64) error
	CheckByID(db *gorm.DB, user *Menu, id int64) error
	Delete(db *gorm.DB, user *entity.Menu, id int64) error
	All(db *gorm.DB, user *[]Menu) error
}

type repository struct {
}

func Newrepository(config *viper.Viper) Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, data *entity.Menu) error {
	return db.Table("menu").Create(&data).Error
}

func (r *repository) Update(db *gorm.DB, data *entity.Menu, id int64) error {
	return db.Table("menu").Updates(&data).Error
}

func (r *repository) CheckByID(db *gorm.DB, data *Menu, id int64) error {
	return db.Table("menu").
		Where("id", id).
		First(&data).
		Scan(data).Error
}

func (r *repository) Delete(db *gorm.DB, data *entity.Menu, id int64) error {
	return db.Table("menu").Delete(&data).Error
}

func (r *repository) All(db *gorm.DB, data *[]Menu) error {
	return db.Table("menu").Find(&data).Error
}
