package test

import (
	"github.com/gookit/slog"
)

type TestRepository interface {
	TestRepository() error
}

type testRepository struct {
	Log *slog.Logger
}

func NewTestRepository(log *slog.Logger) TestRepository {
	return &testRepository{
		Log: log,
	}
}

func (r testRepository) TestRepository() error {
	return nil
}
