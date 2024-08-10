package auth

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CheckUsersByUsername(db *gorm.DB, user *LoginUsers, username string) error
	// UpdateUsers(db *gorm.DB, user *entity.Users, code string) error
	// GetUserByEmail(db *gorm.DB, email string, loginBy string) (*UserData, error)
	// RegisterUserBySso(db *gorm.DB, req RegisterUserRequest) (*UserData, error)
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
