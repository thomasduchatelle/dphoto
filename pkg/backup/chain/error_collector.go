package chain

import (
	"github.com/pkg/errors"
	"sync"
)

type ErrorObserver func(error)

// NewErrorCollector collects errors that occurred during an async process (is thread-safe)
func NewErrorCollector(observers ...ErrorObserver) ChainableErrorCollector {
	return &errorCollector{
		errorsMutex: sync.Mutex{},
		observers:   observers,
	}
}

type errorCollector struct {
	errors      []error
	errorsMutex sync.Mutex
	observers   []ErrorObserver
}

func (e *errorCollector) OnError(err error) {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	e.errors = append(e.errors, err)
	for _, observer := range e.observers {
		observer(err)
	}
}

func (e *errorCollector) Error() error {
	e.errorsMutex.Lock()
	defer e.errorsMutex.Unlock()

	if len(e.errors) == 0 {
		return nil
	}

	return errors.Wrapf(e.errors[0], "%d error(s) reported before shutdown. First one encountered", len(e.errors))
}
