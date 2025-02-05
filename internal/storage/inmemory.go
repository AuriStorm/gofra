package storage

import (
	"context"
	"errors"
	"gofra/internal/config"
	"sync"
)

var (
	ErrQueueFull        = errors.New("queue_full")
	ErrMaxQueuesReached = errors.New("max_queues_reached")
	ErrTimeoutReached   = errors.New("timeout_reached")
)

/*
	when I was researching, deepseek suggested to use container/list also instad of struct
	but it seems more pythonish way to use this lib and not very popular in go community

	this also highlight me solution with sync.Cond with Wait() and Signal() notification
	but it looks more expensive technique. Doubt also that raw AI output will works correctly
*/

type MessageQueue chan string

type QueueNode struct {
	name string
	msg  chan string
}

type InmemoryQueue struct {
	// AFAIK exclusive lock guarantee the ordering on thread read attempts, so we can use this instead of WaitGroup
	mu          *sync.Mutex
	qNodes      map[string]*QueueNode
	maxQueueCnt int32
	queueSize   int32
}

func NewInmemoryQueue(conf config.StorageConfig) InmemoryQueue {
	return InmemoryQueue{
		qNodes:      make(map[string]*QueueNode),
		mu:          &sync.Mutex{},
		maxQueueCnt: conf.MaxQueueCnt,
		queueSize:   conf.QueueSize,
	}
}

func (iq *InmemoryQueue) Get(ctx context.Context, qName string) (string, error) {
	iq.mu.Lock()
	queue, ok := iq.qNodes[qName]
	iq.mu.Unlock()
	if !ok {
		if err := iq.createQueue(qName); err != nil {
			if errors.Is(err, ErrMaxQueuesReached) {
				return "", ErrMaxQueuesReached
			}
		}

		iq.mu.Lock()
		queue = iq.qNodes[qName]
		iq.mu.Unlock()
	}

	select {
	case msg := <-queue.msg:
		return msg, nil
	case <-ctx.Done():
		return "", ErrTimeoutReached
	}
}

func (iq *InmemoryQueue) Put(ctx context.Context, qName string, message string) error {
	iq.mu.Lock()
	queue, ok := iq.qNodes[qName]
	iq.mu.Unlock()

	if !ok {
		if err := iq.createQueue(qName); err != nil {
			if errors.Is(err, ErrMaxQueuesReached) {
				return ErrMaxQueuesReached
			}
		}

		iq.mu.Lock()
		queue = iq.qNodes[qName]
		iq.mu.Unlock()
	}

	select {
	case queue.msg <- message:
		return nil
	case <-ctx.Done():
		return ErrTimeoutReached
	default:
		return ErrQueueFull
	}
}

func (iq *InmemoryQueue) createQueue(qName string) error {
	iq.mu.Lock()
	defer iq.mu.Unlock()

	if int32(len(iq.qNodes)) >= iq.maxQueueCnt {
		return ErrMaxQueuesReached
	}

	_, ok := iq.qNodes[qName]
	if !ok {
		iq.qNodes[qName] = &QueueNode{
			name: qName,
			msg:  make(chan string, iq.queueSize),
		}
	}
	return nil
}
