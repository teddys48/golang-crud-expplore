package auth

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CheckUsersByUsername(db *gorm.DB, user *LoginUsers, username string) error
	CheckUsersByEmailOrNIP(db *gorm.DB, user *LoginUsers, email_or_nip string) error
}

type authRepository struct {
	// config *viper.Viper
}

func NewAuthRepository(config *viper.Viper) AuthRepository {
	return &authRepository{
		// config: config,
	}
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
		Select("username", "email", "code", "password", "id").
		Where("status", "ACTIVE").
		Where("email", email_or_nip).
		Or("nip", email_or_nip).
		First(&user).
		Scan(user).Error
}
