package config

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func NewLogger() *slog.Logger {
	log := slog.New()
	allLog := handler.MustRotateFile("./logs/log.log", rotatefile.EveryDay, handler.WithLogLevels(slog.AllLevels))
	errLog := handler.MustFileHandler("./logs/error.log", handler.WithLogLevels(slog.DangerLevels))

	slog.PushHandlers(allLog, errLog)

	return log
}
