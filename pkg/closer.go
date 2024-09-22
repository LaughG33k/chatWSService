package pkg

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var C Closer = &closer{CloseFuncs: make([]func(context.Context) error, 0), mu: sync.Mutex{}}

type Closer interface {
	Close(ctx context.Context) error
	Add(func(context.Context) error)
}

type closer struct {
	CloseFuncs []func(context.Context) error
	mu         sync.Mutex
}

func (c *closer) Add(fn func(context.Context) error) {
	c.mu.Lock()
	c.CloseFuncs = append(c.CloseFuncs, fn)
	c.mu.Unlock()
}

func (c *closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	errs := make([]error, len(c.CloseFuncs))
	wg := &sync.WaitGroup{}
	i := &atomic.Int32{}
	done := make(chan struct{}, 1)
	wg.Add(len(c.CloseFuncs))
	i.Store(int32(len(c.CloseFuncs)))

	for _, f := range c.CloseFuncs {

		fn := f
		go func() {
			if err := fn(ctx); err != nil {
				errs = append(errs, err)
			}
			i.Add(-1)
			wg.Done()
		}()

	}
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {

	case <-ctx.Done():
		if i.Load() > 0 {
			return errors.New(fmt.Sprintf("Closing shutdown by context. Not closed services: %d", i.Load()))
		}

	case <-done:
		if len(errs) > 0 {
			return errors.Join(errs...)
		}

	}

	return nil
}
