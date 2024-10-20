package backup

import (
	"github.com/pkg/errors"
	"sync"
)

// NewErrorCollectorObserver collects errors that occurred during the backup process
func NewErrorCollectorObserver() *ErrorCollectorObserver {
	return &ErrorCollectorObserver{
		errorsMutex: sync.Mutex{},
	}
}

type ErrorCollectorObserver struct {
	errors      []error
	errorsMutex sync.Mutex
}

func (e *ErrorCollectorObserver) OnRejectedMedia(found FoundMedia, err error) {
	e.appendError(errors.Wrapf(err, "error in analyser"))
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
