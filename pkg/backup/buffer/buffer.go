package buffer

import "context"

// NewBuffer creates a buffer with a given size.
func NewBuffer[T any](bufferSize int, consumer BufferedConsumer[T]) *Buffer[T] {
	if bufferSize <= 0 {
		bufferSize = 1
	}

	return &Buffer[T]{
		content:  make([]T, 0, bufferSize),
		consumer: consumer,
	}
}

type BufferedConsumer[T any] func(ctx context.Context, buffer []T) error

type Buffer[T any] struct {
	consumer BufferedConsumer[T]
	content  []T
}

func (b *Buffer[T]) Append(ctx context.Context, item T) error {
	b.content = append(b.content, item)

	if len(b.content) >= cap(b.content) {
		err := b.consumer(ctx, b.content)
		b.content = make([]T, 0, cap(b.content))
		return err
	}

	return nil
}

func (b *Buffer[T]) Flush(ctx context.Context) error {
	if len(b.content) == 0 {
		return nil
	}

	return b.consumer(ctx, b.content)
}
