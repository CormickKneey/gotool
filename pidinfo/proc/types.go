package proc

import (
	"errors"
	"fmt"
	"runtime"
)

type procConfig struct {
	basepath string
	contents map[string]string
}

// ProcErr All errors returned should be type ProcErr and include a stack
type ProcErr struct {
	error
	Message string
	Stack   []byte
}

func (f *ProcErr) Error() string {
	return fmt.Sprintf("procreader: %s", f.Message)
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*ProcErr); ok {
		return e
	}

	trace := make([]byte, 4096)
	runtime.Stack(trace, true)

	return &ProcErr{
		error:   err,
		Message: err.Error(),
		Stack:   trace,
	}
}

func newError(format string, args ...interface{}) error {
	return wrapError(errors.New(fmt.Sprintf(format, args...)))
}
