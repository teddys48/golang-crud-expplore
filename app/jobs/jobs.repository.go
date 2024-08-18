package jobs

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Repository interface {
	Create(db *gorm.DB, user *Jobs) error
	Update(db *gorm.DB, user *Jobs, id int64) error
	CheckByID(db *gorm.DB, user *Jobs, id int64) error
	Delete(db *gorm.DB, user *Jobs, id int64) error
	All(db *gorm.DB, user *[]Jobs) error
	GetProjectByID(db *gorm.DB, data *[]Project, id []int64) error
	CheckByProjectID(db *gorm.DB, data *[]Jobs, id int64) error
	Disapprove(db *gorm.DB, data *Jobs, id int64) error
}

type repository struct {
}

func Newrepository(config *viper.Viper) Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, data *Jobs) error {
	return db.Table("jobs").Create(&data).Error
}

func (r *repository) Update(db *gorm.DB, data *Jobs, id int64) error {
	return db.Table("jobs").
		Where("id", id).
		Updates(&data).Error
}

func (r *repository) CheckByID(db *gorm.DB, data *Jobs, id int64) error {
	return db.Table("jobs").
		Where("id", id).
		Where("status", "ACTIVE").
		First(&data).
		Scan(data).Error
}

func (r *repository) Delete(db *gorm.DB, data *Jobs, id int64) error {
	return db.Table("jobs").Delete(&data).Error
}

func (r *repository) All(db *gorm.DB, data *[]Jobs) error {
	return db.Table("jobs").Find(&data).Error
}

func (r *repository) GetProjectByID(db *gorm.DB, data *[]Project, id []int64) error {
	return db.Table("project").Where("id in (?)", id).First(&data).Scan(data).Error
}

func (r *repository) CheckByProjectID(db *gorm.DB, data *[]Jobs, id int64) error {
	return db.Table("jobs").
		Where("project_id", id).
		Where("status", "ACTIVE").
		Find(&data).
		Scan(data).Error
}

func (r *repository) Disapprove(db *gorm.DB, data *Jobs, id int64) error {
	return db.Table("jobs").
		Where("id", id).
		Updates(map[string]interface{}{"approved_by": nil, "approved_on": nil}).Error
}
