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

type StartLink interface {
	// Starts is called first, it should create the channels and start the goroutines ; Next links Starts should also be called.
	Starts(ctx context.Context, collector ChainableErrorCollector) error

	// WaitForCompletion is called after NotifyUpstreamCompleted and should return the error collected
	WaitForCompletion() chan error
}

type Link[Consumed any, Produced any] interface {
	StartLink
	Consumer[Consumed]

	// NotifyUpstreamCompleted is called when the previous link will not call Consume anymore
	NotifyUpstreamCompleted()
}

// MultithreadedLink runs the Operator on as many routines as requested.
type MultithreadedLink[Consumed any, Produced any] struct {
	NumberOfRoutines int
	ConsumerBuilder  func(Consumer[Produced]) Consumer[Consumed]
	Next             Link[Produced, any]
	channel          chan Consumed
}

func (m *MultithreadedLink[Consumed, Produced]) ChainNextLink(next Link[Produced, any]) {
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

type SliceLauncher[Produced any] struct {
	Next     Link[Produced, any]
	Producer func(ctx context.Context) ([]Produced, error)
}

func (s *SliceLauncher[Produced]) Starts(ctx context.Context, collector ChainableErrorCollector) error {
	err := s.Next.Starts(ctx, collector)
	if err != nil {
		return err
	}

	produced, err := s.Producer(ctx)
	if err != nil {
		return err
	}

	for _, product := range produced {
		err = s.Next.Consume(ctx, product)
		if err != nil {
			collector.OnError(err)
			return nil
		}
	}
	s.Next.NotifyUpstreamCompleted()

	return err
}

func (s *SliceLauncher[Produced]) WaitForCompletion() chan error {
	return s.Next.WaitForCompletion()
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
