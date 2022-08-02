package entry

import (
	"runtime"
	"time"

	"github.com/zhangel/go-frame/log/fields"
	"github.com/zhangel/go-frame/log/level"
)

type Entry struct {
	Fields fields.Fields
	Time   time.Time
	Level  level.Level
	Caller *runtime.Frame
	Msg    string
}

type TraceFunctionNameOption struct {
	Enabled bool
}
