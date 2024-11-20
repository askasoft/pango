package goqulog

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
)

type GoquLogger struct {
	Logger log.Logger
	Level  log.Level
}

func NewGoquLogger(logger log.Logger, levels ...log.Level) *GoquLogger {
	gl := &GoquLogger{
		Logger: logger,
		Level:  log.LevelDebug,
	}

	if len(levels) > 0 {
		gl.Level = levels[0]
	}

	return gl
}

func (gl *GoquLogger) Printf(format string, v ...any) {
	lvl := gl.Level

	if gl.Logger.IsLevelEnabled(lvl) {
		le := &log.Event{
			Name:  gl.Logger.GetName(),
			Props: gl.Logger.GetProps(),
			Level: lvl,
			Msg:   fmt.Sprintf(format, v...),
			Time:  time.Now(),
		}
		le.CallerStop("/goqu/", gl.Logger.GetTraceLevel() >= lvl)

		gl.Logger.Write(le)
	}
}
