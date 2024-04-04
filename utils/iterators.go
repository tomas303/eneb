package utils

type IteratorChannel[T any] chan (Result[T, error])

type IteratorFunc[T any] func(channel IteratorChannel[T])

type Iterator[T any] struct {
	current T
	channel IteratorChannel[T]
	link    func()
}

func (iterator *Iterator[T]) Next() bool {
	// fmt.Println("reding")
	result, ok := <-iterator.channel
	// fmt.Println("read", result.Value)
	if ok && result.Err == nil {
		iterator.current = result.Value
	}
	return ok && result.Err == nil
}

func (iterator *Iterator[T]) Current() T {
	return iterator.current
}

func NewIterator[T any](iteratorFunc IteratorFunc[T]) *Iterator[T] {
	var iterator Iterator[T]
	iterator.channel = make(IteratorChannel[T])
	iterator.link = func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "send on closed channel" {
					panic(r)
				}
			}
		}()
		iteratorFunc(iterator.channel)
	}
	go iterator.link()
	return &iterator
}
