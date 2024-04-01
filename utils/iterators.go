package utils

import "sync"

type IteratorChannel[T any] chan (Result[*T, error])

type IteratorFunc[T any] func(channel IteratorChannel[T], accessor *Accessor)

type Iterator[T any] struct {
	accessor Accessor
	current  *T
	channel  IteratorChannel[T]
	link     func()
}

func (iterator *Iterator[T]) Next() bool {
	iterator.accessor.Signal()
	result, ok := <-iterator.channel
	if ok && result.Err != nil {
		iterator.current = result.Value
	}
	return ok && result.Err != nil
}

func (iterator *Iterator[T]) Current() *T {
	return iterator.current
}

func NewIterator[T any](iteratorFunc IteratorFunc[T]) *Iterator[T] {
	var iterator Iterator[T]
	iterator.channel = make(IteratorChannel[T])
	iterator.link = func() {
		iteratorFunc(iterator.channel, &iterator.accessor)
	}
	go iterator.link()
	return &iterator
}

type Accessor struct {
	mutex sync.Mutex
	cond  sync.Cond
	ready bool
}

func (accessor *Accessor) Wait() {
	accessor.mutex.Lock()
	for !accessor.ready {
		accessor.cond.Wait()
	}
	accessor.ready = false
	accessor.mutex.Unlock()
}

func (accessor *Accessor) Signal() {
	accessor.mutex.Lock()
	accessor.ready = true
	accessor.cond.Signal()
	accessor.mutex.Unlock()
}
