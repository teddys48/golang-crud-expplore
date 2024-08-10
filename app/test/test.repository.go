package test

import (
	"github.com/gookit/slog"
	"gorm.io/gorm"
)

type TestRepository interface {
	TestRepository(tx *gorm.DB) error
}

type testRepository struct {
	Log *slog.Logger
}

func NewTestRepository(log *slog.Logger) TestRepository {
	return &testRepository{
		Log: log,
	}
}

func (r testRepository) TestRepository(db *gorm.DB) error {
	return db.Table("auth.users").Error
}
