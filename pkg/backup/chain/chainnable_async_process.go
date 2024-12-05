package chain

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

// Chain is a DSL method to create a chain of processes.
//func Chain[Linked any, After any](rootLink Link[any, Linked], Next Link[Linked, After]) Link[Linked, After] {
//	rootLink.ChainNextLink(Next.(Link[Linked, any]))
//	return Next
//}

//type Chainable[Produced any] interface {
//
//}

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
	NumberOfRoutines int
	ConsumerBuilder  func(Consumer[Produced]) Consumer[Consumed]
	Next             Link[Produced]
	channel          chan Consumed
}

func (m *MultithreadedLink[Consumed, Produced]) ChainNextLink(next Link[Produced]) {
	m.Next = next
}

func (m *MultithreadedLink[Consumed, Produced]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	if m.ConsumerBuilder == nil {
		return errors.New("ConsumerBuilder is not set")
	}
	if m.Next == nil {
		return errors.New("Next is not set")
	}

	m.channel = make(chan Consumed, 255)
	if m.NumberOfRoutines <= 0 {
		m.NumberOfRoutines = 1
	}

	err := m.Next.Starts(ctx, collector)
	if err != nil {
		return err
	}

	consumer := m.ConsumerBuilder(m.Next)

	startsInParallel(ctx, m.NumberOfRoutines, func(ctx context.Context) {
		for consumed := range m.channel {
			err := consumer.Consume(ctx, consumed)
			if err != nil {
				collector.OnError(err)
			}
		}
	}, m.Next.NotifyUpstreamCompleted)

	return nil
}

func (m *MultithreadedLink[Consumed, Produced]) Consume(ctx context.Context, consumed Consumed) error {
	m.channel <- consumed
	return nil
}

func (m *MultithreadedLink[Consumed, Produced]) NotifyUpstreamCompleted() {
	close(m.channel)
}

func (m *MultithreadedLink[Consumed, Produced]) WaitForCompletion() chan error {
	return m.Next.WaitForCompletion()
}

// EnderChainLink runs Operator on the same routines as the previous link, and return errors from the ChainableErrorCollector
type EnderChainLink[Consumed any] struct {
	done      chan error
	Operator  func(ctx context.Context, consumed Consumed) error
	collector ChainableErrorCollector
}

func (l *EnderChainLink[Consumed]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	l.done = make(chan error)
	l.collector = collector
	return nil
}

func (l *EnderChainLink[Consumed]) Consume(ctx context.Context, produced Consumed) error {
	if l.Operator != nil {
		return l.Operator(ctx, produced)
	}

	return nil
}

func (l *EnderChainLink[Consumed]) NotifyUpstreamCompleted() {
	if err := l.collector.Error(); err != nil {
		l.done <- err
	}
	close(l.done)
}

func (l *EnderChainLink[Consumed]) WaitForCompletion() chan error {
	return l.done
}

// SingleLauncher launch the chain process by consume one and only one element.
type SingleLauncher[Consumed any, Produced any] struct {
	Next      Link[Produced]
	Function  func(ctx context.Context, consumed Consumed) ([]Produced, error)
	collector ChainableErrorCollector
}

func (s *SingleLauncher[Consumed, Produced]) Consume(ctx context.Context, consumed Consumed) error {
	defer s.Next.NotifyUpstreamCompleted()

	products, err := s.Function(ctx, consumed)
	if err != nil {
		return err
	}

	for _, product := range products {
		err := s.Next.Consume(ctx, product)
		if err != nil {
			s.collector.OnError(err)
			return nil
		}
	}

	return nil
}

func (s *SingleLauncher[Consumed, Produced]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	s.collector = collector
	return s.Next.Starts(ctx, collector)
}

func (s *SingleLauncher[Consumed, Produced]) WaitForCompletion() chan error {
	return s.Next.WaitForCompletion()
}

// Process combine Consume and WaitForCompletion to simplify consumption.
func (s *SingleLauncher[Consumed, Produced]) Process(ctx context.Context, consumed Consumed) chan error {
	err := s.Consume(ctx, consumed)
	if err != nil {
		errChan := make(chan error, 1)
		errChan <- err
		return errChan
	}

	return s.WaitForCompletion()
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
