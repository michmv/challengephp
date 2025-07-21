package lib

import (
	"fmt"
	"runtime"
	"strings"
)

type Error interface {
	Add(string) Error
	Tap() Error
	Error() string
}

type LogError struct {
	err      error
	messages []string
	trace    []string
	stack    []byte
}

func Err(err error) *LogError {
	return &LogError{
		err:   err,
		trace: []string{trace()},
	}
}

func (it *LogError) Add(msg string) Error {
	it.trace = append(it.trace, trace())
	it.err = fmt.Errorf("%s: %w", msg, it.err)
	return it
}

func (it *LogError) Tap() Error {
	it.trace = append(it.trace, trace())
	return it
}

func (it *LogError) Error() string {
	result := ""
	if len(it.trace) > 0 {
		result = "\n\t" + strings.Join(it.trace, "\n\t")
	}
	return it.err.Error() + result
}

func trace() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "it was not possible to recover the information"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

type PanicError struct {
	*LogError
}

func ErrPanic(err error, stack []byte) *PanicError {
	e := PanicError{
		LogError: Err(err),
	}
	e.trace = strings.Split(string(stack), "\n")
	return &e
}
