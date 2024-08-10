package config

import (
	"fmt"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func NewLogger() {
	// log := slog.New()
	// myTemplate := "[{{datetime}}] [{{level}}] {{message}}"
	// a := slog.new

	// x := slog.NewStdLogger(func(sl *slog.SugaredLogger) {
	// 	sl.Formatter = a
	// })
	allLog := handler.MustFileHandler("./logs/log", handler.WithRotateTime(rotatefile.EveryDay), handler.WithLogLevels(slog.AllLevels), func(c *handler.Config) {
		// c.RotateWriter()
		c.BackupTime = 1
		c.RotateMode = rotatefile.ModeCreate
		// c.CreateHandler()
		c.RotateMode = rotatefile.RotateMode(1)
		c.RotateMode = rotatefile.ModeCreate
		c.RenameFunc = func(filepath string, rotateNum uint) string {
			fmt.Println(filepath)
			newName := fmt.Sprintf("%s-%s.log", filepath, time.Now().Format("2006-01-02"))
			if rotateNum > 0 {
				fmt.Println("qwqwqwasasas")
				newName = fmt.Sprintf("%s.%d", newName, rotateNum)
			}
			return newName
		}
		// c.TimeClock = rotatefile.DefaultTimeClockFn
		// c.CreateWriter()
	})
	errLog := handler.MustFileHandler("./logs/error.log", handler.WithLogLevels(slog.DangerLevels))
	slog.PushHandlers(allLog, errLog)
	// slog.PushHandlers(allLog, errLog)
}
