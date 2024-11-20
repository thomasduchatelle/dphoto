package chain

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/backup/buffer"
)

// BufferLink use a channel to collect items and release them as soon as the BufferCapacity is reached. Next is running on a single thread.
type BufferLink[Consumed any] struct {
	BufferCapacity int
	Next           Link[[]Consumed]
	channel        chan Consumed
	buffer         *buffer.Buffer[Consumed]
}

func (l *BufferLink[Consumed]) Consume(ctx context.Context, consumed Consumed) error {
	l.channel <- consumed
	return nil
}

func (l *BufferLink[Consumed]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	l.channel = make(chan Consumed, 255)
	l.buffer = buffer.NewBuffer(l.BufferCapacity, l.Next.Consume)

	go func() {
		defer l.Next.NotifyUpstreamCompleted()

		for {
			select {
			case consumed, more := <-l.channel:
				if more {
					err := l.buffer.Append(ctx, consumed)
					if err != nil {
						collector.OnError(err)
					}
				} else {
					err := l.buffer.Flush(ctx)
					if err != nil {
						collector.OnError(err)
					}
					return
				}
			}
		}
	}()

	return l.Next.Starts(ctx, collector)
}

func (l *BufferLink[Consumed]) WaitForCompletion() chan error {
	return l.Next.WaitForCompletion()
}

func (l *BufferLink[Consumed]) NotifyUpstreamCompleted() {
	close(l.channel)
}

type ReBufferLink[Consumed any] struct {
	BufferLink[Consumed]
}

func (l *ReBufferLink[Consumed]) Consume(ctx context.Context, buf []Consumed) error {
	for _, item := range buf {
		err := l.BufferLink.Consume(ctx, item)
		if err != nil {
			return err
		}
	}
	return nil
}
