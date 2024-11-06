package backup

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

type errorObserver func(error)

// TODO Privatise the error collector, drop the interface, drop the "observer" implementation (OnRejectedMedia).

type IErrorCollectorObserver interface {
	RejectedMediaObserver

	appendError(err error)
	hasAnyErrors() int
	Errors() []error
}

// newErrorCollector collects errors that occurred during an async process (is thread-safe)
func newErrorCollector() *errorCollector {
	return &errorCollector{
		errorsMutex: sync.Mutex{},
	}
}

type errorCollector struct {
	errors         []error
	errorObservers []errorObserver
	errorsMutex    sync.Mutex
}

func (e *errorCollector) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	e.appendError(errors.Wrapf(cause, "error in analyser"))
	return nil
}

func (e *errorCollector) appendError(err error) {
	if err == nil {
		return
	}

	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	e.errors = append(e.errors, err)
	for _, observer := range e.errorObservers {
		observer(err)
	}
}

func (e *errorCollector) hasAnyErrors() int {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	return len(e.errors)
}

func (e *errorCollector) collectError() error {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	if len(e.errors) == 0 {
		return nil
	}

	return errors.Wrapf(e.errors[0], "%d error(s) reported before shutdown. First one encountered", len(e.errors))
}

// Errors TODO - remove this method and use collectError instead
func (e *errorCollector) Errors() []error {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	errs := make([]error, len(e.errors), len(e.errors))
	copy(errs, e.errors)

	return errs
}

func (e *errorCollector) registerErrorObserver(observer errorObserver) {
	e.errorObservers = append(e.errorObservers, observer)
}
