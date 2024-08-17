package corporation

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Repository interface {
	Create(db *gorm.DB, user *Corporation) error
	Update(db *gorm.DB, user *Corporation, id int64) error
	CheckByID(db *gorm.DB, user *Corporation, id int64) error
	Delete(db *gorm.DB, user *Corporation, id int64) error
	All(db *gorm.DB, user *[]Corporation) error
}

type repository struct {
}

func Newrepository(config *viper.Viper) Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, data *Corporation) error {
	return db.Table("corporation").Create(&data).Error
}

func (r *repository) Update(db *gorm.DB, data *Corporation, id int64) error {
	return db.Table("corporation").
		Where("id", id).
		Updates(&data).Error
}

func (r *repository) CheckByID(db *gorm.DB, data *Corporation, id int64) error {
	return db.Table("corporation").
		Where("id", id).
		Where("status", "ACTIVE").
		First(&data).
		Scan(data).Error
}

func (r *repository) Delete(db *gorm.DB, data *Corporation, id int64) error {
	return db.Table("corporation").Delete(&data).Error
}

func (r *repository) All(db *gorm.DB, data *[]Corporation) error {
	return db.Table("corporation").Find(&data).Error
}
