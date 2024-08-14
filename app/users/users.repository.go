package users

import (
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/entity"
	"gorm.io/gorm"
)

type Repository interface {
	CheckUsersByUsername(db *gorm.DB, user *Users, username string) error
	Create(db *gorm.DB, user *entity.Users) error
	Update(db *gorm.DB, user *entity.Users, id int64) error
	CheckUsersByID(db *gorm.DB, user *Users, id int64) error
	Delete(db *gorm.DB, user *entity.Users, id int64) error
	All(db *gorm.DB, user *[]Users) error
}

type repository struct {
}

func Newrepository(config *viper.Viper) Repository {
	return &repository{}
}

func (r *repository) CheckUsersByUsername(db *gorm.DB, user *Users, username string) error {
	return db.Table("users").
		Where("username", username).
		Where("status", "ACTIVE").
		First(&user).
		Scan(user).Error
}

func (r *repository) Create(db *gorm.DB, user *entity.Users) error {
	return db.Table("users").
		Create(&user).Error
}

func (r *repository) Update(db *gorm.DB, user *entity.Users, id int64) error {
	return db.Table("users").
		Where("id", id).
		Updates(&user).Error
}

func (r *repository) CheckUsersByID(db *gorm.DB, user *Users, id int64) error {
	return db.Table("users").
		Where("id", id).
		Where("status", "ACTIVE").
		First(&user).
		Scan(user).Error
}

func (r *repository) Delete(db *gorm.DB, user *entity.Users, id int64) error {
	return db.Table("users").
		Where("id", id).
		Delete(&user).Error
}

func (r *repository) All(db *gorm.DB, user *[]Users) error {
	return db.Table("users").
		Find(&user).
		Scan(user).Error
}
