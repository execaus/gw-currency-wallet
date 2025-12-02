package pkg

import (
	"context"
	"sync"
	"time"
)

type UpdateFn[T any] func(ctx context.Context) (T, error)

const (
	DefaultCacherTTL = 5 * time.Second
)

type Cacher[T any] struct {
	data     T
	updateFn UpdateFn[T]
	ttl      time.Duration
	lastSync time.Time
	mu       sync.Mutex
	cond     *sync.Cond
}

func NewCacher[T any](ctx context.Context, updateFn UpdateFn[T], ttl time.Duration) (*Cacher[T], error) {
	data, err := updateFn(ctx)
	if err != nil {
		return nil, err
	}
	return &Cacher[T]{
		data:     data,
		updateFn: updateFn,
		ttl:      ttl,
		lastSync: time.Now(),
	}, nil
}

func (c *Cacher[T]) GetData(ctx context.Context) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.lastSync) < c.ttl {
		return c.data, nil
	}

	if c.cond != nil {
		c.cond.Wait()
		return c.data, nil
	}

	c.cond = sync.NewCond(&c.mu)
	cond := c.cond
	c.mu.Unlock()

	data, err := c.updateFn(ctx)

	c.mu.Lock()
	c.cond = nil
	cond.Broadcast()

	if err != nil {
		return c.data, err
	}

	c.data = data
	c.lastSync = time.Now()
	return c.data, nil
}

func (c *Cacher[T]) ForceSync(ctx context.Context) (T, error) {
	c.mu.Lock()
	if c.cond != nil {
		c.cond.Wait()
		data := c.data
		c.mu.Unlock()
		return data, nil
	}

	c.cond = sync.NewCond(&c.mu)
	cond := c.cond
	c.mu.Unlock()

	data, err := c.updateFn(ctx)

	c.mu.Lock()
	c.cond = nil
	cond.Broadcast()

	if err == nil {
		c.data = data
		c.lastSync = time.Now()
	}
	result := c.data
	c.mu.Unlock()

	return result, err
}
