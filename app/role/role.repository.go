package role

import (
	"github.com/spf13/viper"
	"github.com/teddys48/kmpro/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(db *gorm.DB, user *entity.Role) error
	Update(db *gorm.DB, user *entity.Role, id int64) error
	CheckByID(db *gorm.DB, user *Role, id int64) error
	Delete(db *gorm.DB, user *entity.Role, id int64) error
	All(db *gorm.DB, user *[]Role) error
	GetRoleDetailData(db *gorm.DB, data *[]RoleDetailData, role_id int) error
	InsertUpdateRole(db *gorm.DB, dataRole *entity.Role, dataRoleDetail *[]entity.RoleDetail) error
	UpdateUpdateRole(db *gorm.DB, dataRole *entity.Role, dataRoleDetail *[]entity.RoleDetail, roleID int64) error
}

type repository struct {
}

func Newrepository(config *viper.Viper) Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, data *entity.Role) error {
	return db.Table("role").Create(&data).Error
}

func (r *repository) Update(db *gorm.DB, data *entity.Role, id int64) error {
	return db.Table("role").
		Where("id", id).
		Updates(&data).Error
}

func (r *repository) CheckByID(db *gorm.DB, data *Role, id int64) error {
	return db.Table("role").
		Where("id", id).
		Where("status", "ACTIVE").
		First(&data).
		Scan(data).Error
}

func (r *repository) Delete(db *gorm.DB, data *entity.Role, id int64) error {
	return db.Table("role").Delete(&data).Error
}

func (r *repository) All(db *gorm.DB, data *[]Role) error {
	return db.Table("role").Find(&data).Error
}

func (r *repository) GetRoleDetailData(db *gorm.DB, data *[]RoleDetailData, role_id int) error {
	return db.Table("role_detail").
		Select(
			"menu.name as name",
			"menu_path_url as path",
			"menu.sort as sort",
			"role_detail.action as action").
		Joins("left join menu on menu.id = role_detail.id").
		Find(&data).
		Scan(data).Error
}

func (r *repository) InsertUpdateRole(db *gorm.DB, dataRole *entity.Role, dataRoleDetail *[]entity.RoleDetail) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Table("role").
		Create(&dataRole).Error; err != nil {
		tx.Rollback()
		return err
	}

	roleID := dataRole.ID

	for i := range *dataRoleDetail {
		(*dataRoleDetail)[i].RoleID = roleID
	}

	if err := tx.Table("role_detail").
		Create(&dataRoleDetail).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *repository) UpdateUpdateRole(db *gorm.DB, dataRole *entity.Role, dataRoleDetail *[]entity.RoleDetail, roleID int64) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Table("role").
		Where("id", roleID).
		Updates(&dataRole).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("role_detail").
		Where("role_id", roleID).
		Create(&dataRoleDetail).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
