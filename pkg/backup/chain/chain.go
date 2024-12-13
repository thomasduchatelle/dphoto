package chain

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

type ChainableErrorCollector interface {
	OnError(err error)
	Error() error
}

type Consumer[Consumed any] interface {
	Consume(ctx context.Context, consumed Consumed) error
}

type ConsumerFunc[Consumed any] func(ctx context.Context, consumed Consumed) error

func (c ConsumerFunc[Consumed]) Consume(ctx context.Context, consumed Consumed) error {
	return c(ctx, consumed)
}

type StartLink[Consumed any] interface {
	Consumer[Consumed]

	// Starts is called first, it should create the channels and start the goroutines ; Next links Starts should also be called.
	Starts(ctx context.Context, collector ChainableErrorCollector) error

	// WaitForCompletion is called after NotifyUpstreamCompleted and should return the error collected
	WaitForCompletion() chan error
}

type Link[Consumed any] interface {
	StartLink[Consumed]

	// NotifyUpstreamCompleted is called when the previous link will not call Consume anymore
	NotifyUpstreamCompleted()
}

// MultithreadedLink runs the Operator on as many routines as requested.
type MultithreadedLink[Consumed any, Produced any] struct {
	NumberOfRoutines int                                         // NumberOfRoutines is the number of routines on which the ConsumerBuilder returned method will be called. Default is 1.
	ConsumerBuilder  func(Consumer[Produced]) Consumer[Consumed] // ConsumerBuilder is the factory function to build the consumer that transforms Consumed into Produced. Use PassThrough if no transformation is needed.
	Cancellable      bool                                        // Cancellable is true if the cancelled context should stop the routine. Default is false.
	ChannelSize      int                                         // ChannelSize is defaulted to 255
	Next             Link[Produced]                              // Next will receive the product of the ConsumerBuilder returned method. It is mandatory to have one, use EndOfTheChain to end the chain.
	channel          chan Consumed
}

func (l *MultithreadedLink[Consumed, Produced]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	if l.ConsumerBuilder == nil {
		return errors.New("MultithreadedLink.ConsumerBuilder is not set")
	}
	if l.Next == nil {
		return errors.New("MultithreadedLink.Next is not set")
	}

	if l.ChannelSize <= 0 {
		l.ChannelSize = 255
	}
	if l.NumberOfRoutines <= 0 {
		l.NumberOfRoutines = 1
	}
	l.channel = make(chan Consumed, l.ChannelSize)

	err := l.Next.Starts(ctx, collector)
	if err != nil {
		return err
	}

	consumer := l.ConsumerBuilder(l.Next)

	var routine func(ctx context.Context)
	if l.Cancellable {
		routine = l.multithreadedLinkCancellableRoutine(consumer, collector)
	} else {
		routine = l.multithreadedLinkDefaultRoutine(consumer, collector)
	}
	startsInParallel(ctx, l.NumberOfRoutines, routine, l.Next.NotifyUpstreamCompleted)

	return nil
}

// multithreadedLinkDefaultRoutine is consuming messages in the channel until the channel is closed. Then the routine terminates.
func (l *MultithreadedLink[Consumed, Produced]) multithreadedLinkDefaultRoutine(consumer Consumer[Consumed], collector ChainableErrorCollector) func(ctx context.Context) {
	return func(ctx context.Context) {
		for consumed := range l.channel {
			err := consumer.Consume(ctx, consumed)
			if err != nil {
				collector.OnError(err)
			}
		}
	}
}

// multithreadedLinkCancellableRoutine is consuming messages until the context is cancelled all or the channel is closed. Then the routine terminates.
func (l *MultithreadedLink[Consumed, Produced]) multithreadedLinkCancellableRoutine(consumer Consumer[Consumed], collector ChainableErrorCollector) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case consumed, more := <-l.channel:
				if !more {
					return
				}

				err := consumer.Consume(ctx, consumed)
				if err != nil {
					collector.OnError(err)
				}
			}
		}
	}
}

func (l *MultithreadedLink[Consumed, Produced]) Consume(ctx context.Context, consumed Consumed) error {
	if l.Cancellable {
		addsToChannelIfContextNotCancelled(ctx, l.channel, consumed)

	} else {
		blockingAddsToChannel(ctx, l.channel, consumed)
	}
	return nil
}

func blockingAddsToChannel[Consumed any](ctx context.Context, channel chan Consumed, consumed Consumed) {
	channel <- consumed
}

func addsToChannelIfContextNotCancelled[Consumed any](ctx context.Context, channel chan Consumed, consumed Consumed) {
	select {
	case <-ctx.Done():
	case channel <- consumed:
	}
}

func (l *MultithreadedLink[Consumed, Produced]) NotifyUpstreamCompleted() {
	close(l.channel)
}

func (l *MultithreadedLink[Consumed, Produced]) WaitForCompletion() chan error {
	return l.Next.WaitForCompletion()
}

func EndOfTheChain[Consumed any](consumers ...ConsumerFunc[Consumed]) *EndLink[Consumed] {
	return &EndLink[Consumed]{
		Consumers: consumers,
	}
}

// EndLink runs Operator on the same routines as the previous link, and return errors from the ChainableErrorCollector
type EndLink[Consumed any] struct {
	done      chan error
	Consumers []ConsumerFunc[Consumed]
	collector ChainableErrorCollector
}

func (l *EndLink[Consumed]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	l.done = make(chan error, 1)
	l.collector = collector
	return nil
}

func (l *EndLink[Consumed]) Consume(ctx context.Context, produced Consumed) error {
	for _, consumer := range l.Consumers {
		err := consumer(ctx, produced)
		if err != nil {
			l.collector.OnError(err)
		}
	}

	return nil
}

func (l *EndLink[Consumed]) NotifyUpstreamCompleted() {
	if err := l.collector.Error(); err != nil {
		l.done <- err
	}
	close(l.done)
}

func (l *EndLink[Consumed]) WaitForCompletion() chan error {
	return l.done
}

// SingleLauncher launch the chain process by consume one and only one element.
type SingleLauncher[Consumed any, Produced any] struct {
	Next      Link[Produced]
	Function  func(ctx context.Context, consumed Consumed) ([]Produced, error)
	ctx       context.Context // ctx used to process the chain is the one used to START the chain, not the one received in Process. This is to keep the chain behaviour consistent no matter the threads.
	collector ChainableErrorCollector
}

func (s *SingleLauncher[Consumed, Produced]) Consume(ctx context.Context, consumed Consumed) error {
	defer s.Next.NotifyUpstreamCompleted()

	products, err := s.Function(ctx, consumed)
	if err != nil {
		return err
	}

	for _, product := range products {
		err = s.Next.Consume(ctx, product)
		if err != nil {
			s.collector.OnError(err)
			return nil
		}
	}

	return nil
}

func (s *SingleLauncher[Consumed, Produced]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	s.ctx = ctx
	s.collector = collector
	return s.Next.Starts(ctx, collector)
}

func (s *SingleLauncher[Consumed, Produced]) WaitForCompletion() chan error {
	return s.Next.WaitForCompletion()
}

// Process combine Consume and WaitForCompletion to simplify consumption.
func (s *SingleLauncher[Consumed, Produced]) Process(ctx context.Context, consumed Consumed) chan error {
	err := s.Consume(s.ctx, consumed)
	if err != nil {
		s.collector.OnError(err)
	}

	return s.WaitForCompletion()
}

type CloserFunc func()

type CloseWrapperLink[Consumed any] struct {
	CloserFuncs []CloserFunc
	Next        Link[Consumed]
}

func (w *CloseWrapperLink[Consumed]) Consume(ctx context.Context, consumed Consumed) error {
	return w.Next.Consume(ctx, consumed)
}

func (w *CloseWrapperLink[Consumed]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	return w.Next.Starts(ctx, collector)
}

func (w *CloseWrapperLink[Consumed]) WaitForCompletion() chan error {
	return w.Next.WaitForCompletion()
}

func (w *CloseWrapperLink[Consumed]) NotifyUpstreamCompleted() {
	for _, f := range w.CloserFuncs {
		f()
	}
	w.Next.NotifyUpstreamCompleted()
}

// PassThrough is a ConsumerBuilder that forwards the .
func PassThrough[Consumed any]() func(Consumer[Consumed]) Consumer[Consumed] {
	return func(c Consumer[Consumed]) Consumer[Consumed] {
		return c
	}
}

func startsInParallel(ctx context.Context, parallel int, consume func(ctx context.Context), closeChannel func()) {
	group := sync.WaitGroup{}

	group.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer group.Done()

			consume(ctx)
		}()
	}

	go func() {
		group.Wait()
		closeChannel()
	}()
}
