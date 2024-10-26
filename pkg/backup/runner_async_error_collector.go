package backup

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

type IErrorCollectorObserver interface {
	RejectedMediaObserver

	appendError(err error)
	hasAnyErrors() int
	Errors() []error
}

// NewErrorCollectorObserver collects errors that occurred during the backup process
func NewErrorCollectorObserver() IErrorCollectorObserver {
	return &ErrorCollectorObserver{
		errorsMutex: sync.Mutex{},
	}
}

type ErrorCollectorObserver struct {
	errors      []error
	errorsMutex sync.Mutex
}

func (e *ErrorCollectorObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	e.appendError(errors.Wrapf(cause, "error in analyser"))
	return nil
}

func (e *ErrorCollectorObserver) appendError(err error) {
	if err == nil {
		return
	}

	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	e.errors = append(e.errors, err)
}

func (e *ErrorCollectorObserver) hasAnyErrors() int {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	return len(e.errors)
}

func (e *ErrorCollectorObserver) Errors() []error {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	errs := make([]error, len(e.errors), len(e.errors))
	copy(errs, e.errors)

	return errs
}
