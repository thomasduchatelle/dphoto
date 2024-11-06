package backup

import (
	"context"
	"slices"
)

type flushableCollector struct {
	flushables []flushable
}

func (f *flushableCollector) append(flushable flushable) {
	f.flushables = append(f.flushables, flushable)
}

func (f *flushableCollector) Flush(ctx context.Context) error {
	for _, flushable := range slices.Backward(f.flushables) {
		if err := flushable.Flush(ctx); err != nil {
			return err
		}
	}

	return nil
}

type flushable interface {
	Flush(ctx context.Context) error
}
