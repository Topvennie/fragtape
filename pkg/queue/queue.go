// Package queue provides queue implementations
package queue

import (
	"context"
	"errors"
	"time"
)

var ErrEmpty = errors.New("queue is empty")

type DequeueOption struct {
	timeout time.Duration
}

type Queue[T any] interface {
	Enqueue(context.Context, T) error
	Dequeue(context.Context, ...DequeueOption) (T, error)
}
