package utils

import (
	"fmt"
)

type IteratorChannel[T any] chan (Result[T, error])

type IteratorFunc[T any] func(channel IteratorChannel[T])

type Iterator[T any] struct {
	current T
	chiter  IteratorChannel[T]
	chcmd   chan bool
	itfunc  func()
}

func (iterator *Iterator[T]) Next() bool {
	// fmt.Println("reding")
	iterator.chcmd <- true
	result, ok := <-iterator.chiter
	// fmt.Println("read", result.Value)
	if ok && result.Err == nil {
		iterator.current = result.Value
	}
	return ok && result.Err == nil
}

func (iterator *Iterator[T]) Current() T {
	return iterator.current
}

func (iterator *Iterator[T]) Close() {
	iterator.chcmd <- false
	close(iterator.chcmd)
	close(iterator.chiter)
}

func NewIterator[T any](iteratorFunc IteratorFunc[T]) *Iterator[T] {
	var iterator Iterator[T]
	iterator.chiter = make(IteratorChannel[T])
	iterator.chcmd = make(chan bool)
	iterator.itfunc = func() {
		for {
			cmd, ok := <-iterator.chcmd
			if cmd && ok {
				iteratorFunc(iterator.chiter)
			} else {
				break
			}
		}
	}
	go iterator.itfunc()
	return &iterator
}

type Iterator2EofErr struct {
	message string
}

func (e Iterator2EofErr) Error() string {
	return e.message
}

var IterEOF = Iterator2EofErr{message: "eof"}

type Iterator2Func[T any] func() (T, bool)
type Iterator2[T any] struct {
	current   T
	itfunc    Iterator2Func[T]
	lastError error
}

func (iterator *Iterator2[T]) Next() bool {
	if iterator.lastError != nil {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("Recovered from panic:", r)
			iterator.lastError = fmt.Errorf("panic: %v", r)
		}
	}()
	var ok bool
	iterator.current, ok = iterator.itfunc()
	if !ok {
		iterator.lastError = IterEOF
	}
	return ok
}

func (iterator *Iterator2[T]) Current() T {
	return iterator.current
}

func (iterator *Iterator2[T]) Error() error {
	return iterator.lastError
}

func NewIterator2[T any](iteratorFunc Iterator2Func[T]) *Iterator2[T] {
	var iterator Iterator2[T]
	iterator.itfunc = iteratorFunc
	return &iterator
}
