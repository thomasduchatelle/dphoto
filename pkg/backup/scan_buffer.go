package backup

import "context"

func bufferAnalysedMedia(ctx context.Context, size int, consumer bufferedConsumer[*AnalysedMedia], register func(flushable flushable)) *bufferAnalysedMediaObserverAdapter {
	adapter := &bufferAnalysedMediaObserverAdapter{
		Buffer: &buffer[*AnalysedMedia]{
			consumer: consumer,
			content:  make([]*AnalysedMedia, 0, size),
		},
	}
	register(adapter.Buffer)

	return adapter
}

type bufferAnalysedMediaObserverAdapter struct {
	Buffer *buffer[*AnalysedMedia]
}

func (b *bufferAnalysedMediaObserverAdapter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	return b.Buffer.Append(ctx, media)
}

type bufferedConsumer[T any] interface {
	callback(ctx context.Context, buffer []T) error
}

type buffer[T any] struct {
	consumer bufferedConsumer[T]
	content  []T
}

func (b *buffer[T]) Append(ctx context.Context, item T) error {
	b.content = append(b.content, item)

	if len(b.content) >= cap(b.content) {
		err := b.consumer.callback(ctx, b.content)
		b.content = make([]T, 0, cap(b.content))
		return err
	}

	return nil
}

func (b *buffer[T]) flush(ctx context.Context) error {
	if len(b.content) == 0 {
		return nil
	}

	return b.consumer.callback(ctx, b.content)
}
