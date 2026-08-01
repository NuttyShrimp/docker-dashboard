package utils

import (
	stderrors "errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

func Error(err error) error {
	if err == nil {
		return nil
	}

	if hasStackTrace(err) {
		return err
	}

	return pkgerrors.WithStack(err)
}

func Errorf(format string, args ...any) error {
	return Error(fmt.Errorf(format, args...))
}

type stackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

func hasStackTrace(err error) bool {
	if err == nil {
		return false
	}

	var st stackTracer
	return stderrors.As(err, &st)
}
