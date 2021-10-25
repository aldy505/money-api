package decorator

import (
	"errors"
	"runtime"
	"strconv"
)

type Frame struct {
	Func string
	Path string
	Line int
}

// Decorates a normal error and fill it with a stack trace.
// A modification of tracerr package.
func Err(err error) error {
	if err == nil {
		return nil
	}

	traced := trace(err, 2)

	var stack string
	for _, v := range traced {
		stack += "\n"
		stack += v.Func
		stack += " "
		stack += v.Path
		stack += ":"
		stack += strconv.Itoa(v.Line)
	}
	return errors.New(err.Error() + "\n" + stack)
}

func trace(err error, skip int) []Frame {
	frames := make([]Frame, 0, 20)
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
	}

	return frames
}
