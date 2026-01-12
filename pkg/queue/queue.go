// Package queue provides queue implementations
package queue

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEmpty = errors.New("queue is empty")
	ErrNoop  = errors.New("no operation")
)

type Receipt struct {
	Raw string
}

type DequeueOption struct {
	timeout time.Duration
}

type PeekOption struct {
	Start int
	End   int
	All   bool
}

// Queue is a FIFO structure
type Queue[T any] interface {
	// Enqueue adds a new items to the queue
	Enqueue(context.Context, T) error
	// Dequeue removes and returns the item that's been the longest in the queue
	Dequeue(context.Context, ...DequeueOption) (T, Receipt, error)
}

// ReliableQueue moves items from the queue to a tmp queue when using Dequeue
// Use the Complete function to remove it from the temporarily queue
type ReliableQueue[T any] interface {
	// The dequeue function will now move items to a temp queue
	Queue[T]
	// Complete marks an item as complete and removes it from the temporarily queue
	Complete(context.Context, Receipt) error
	// RequeueAll moves all items from the tmp queue to the main queue
	// It adds in order that they've been added to the tmp queue
	RequeueAll(context.Context) error
}
