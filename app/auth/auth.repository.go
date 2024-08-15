package auth

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CheckUsersByUsername(db *gorm.DB, user *LoginUsers, username string) error
	CheckUsersByEmailOrNIP(db *gorm.DB, user *LoginUsers, email_or_nip string) error
	GetRoleDetailData(db *gorm.DB, data *[]Menu, role_id int) error
	CheckUsersByEmailOrNIP2(db *gorm.DB, user *UsersData, email_or_nip string) error
}

type authRepository struct {
	// config *viper.Viper
}

func NewAuthRepository(config *viper.Viper) AuthRepository {
	return &authRepository{
		// config: config,
	}
}

func (r authRepository) GetRoleDetailData(db *gorm.DB, data *[]Menu, role_id int) error {
	return db.Table("menu").
		Select(
			"menu.name as name",
			"menu.path_url as path",
			"menu.sort as sort",
			// "role_detail.action as action",
		).
		Joins("left join role_detail on role_detail.menu_id = menu.id").
		Where("role_detail.role_id", role_id).
		// First(&data).
		Scan(data).Error
}

func (r authRepository) CheckUsersByUsername(db *gorm.DB, user *LoginUsers, username string) error {
	return db.Table("users").
		Select("username", "email", "code", "password", "id").
		Where("username", username).
		Where("status", "ACTIVE").
		First(&user).
		Scan(user).Error
}

func (r authRepository) CheckUsersByEmailOrNIP(db *gorm.DB, user *LoginUsers, email_or_nip string) error {
	return db.Table("users").
		Select("username", "email", "code", "password", "id", "role_id").
		Where("status", "ACTIVE").
		Where("email", email_or_nip).
		Or("nip", email_or_nip).
		First(&user).
		Scan(user).Error
}

func (r authRepository) CheckUsersByEmailOrNIP2(db *gorm.DB, user *UsersData, email_or_nip string) error {
	return db.Table("users").
		Select(
			"users.email as email",
			"nip as nip",
			"nrk as nrk",
			"instance as instansi",
			"position as jabatan",
			"full_name as fullname",
			"role.name as role",
			"corporation_id as perusahaan_id",
		).
		Joins("left join role on users.role_id = role.id").
		Where("users.status", "ACTIVE").
		Where("role.status", "ACTIVE").
		Where("email", email_or_nip).
		Or("nip", email_or_nip).
		// First(&user).
		Scan(user).Error
}
