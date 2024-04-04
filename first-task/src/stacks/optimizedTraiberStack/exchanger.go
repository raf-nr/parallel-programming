package optimizedTraiberStack

import (
	"errors"
	"sync/atomic"
)

type status int

type exchanger[T any] struct {
	elem atomic.Value
}

const (
	empty   status = 0
	waiting status = 1
	busy    status = 2
)

type exchangerElem[T any] struct {
	value *T
	state status
}

func freshExchanger[T any]() exchanger[T] {
	fresh := exchanger[T]{}
	fresh.elem.Store(exchangerElem[T]{value: nil, state: empty})
	return fresh
}

func (e *exchanger[T]) exchange(value *T, replays int) (*T, error) {

	for i := 0; i < replays; i++ {
		currentElem := e.elem.Load().(exchangerElem[T])

		switch currentElem.state {
		case empty:
			oldItem := exchangerElem[T]{nil, empty}
			newItem := exchangerElem[T]{value, waiting}
			if e.elem.CompareAndSwap(oldItem, newItem) {
				for j := i; j < replays; j++ {
					currentElem := e.elem.Load().(exchangerElem[T])
					if currentElem.state == busy {
						newItem := exchangerElem[T]{nil, empty}
						e.elem.Store(newItem)
						return currentElem.value, nil
					}
				}
				e.elem.Store(oldItem)
				return new(T), errors.New("The allowed number of iterations has been exceeded.")
			}
		case waiting:
			oldItem := exchangerElem[T]{currentElem.value, waiting}
			newItem := exchangerElem[T]{value, busy}
			if (oldItem.value == nil) != (value == nil) {
				if e.elem.CompareAndSwap(oldItem, newItem) {
					return currentElem.value, nil
				}
			}
		case busy:
			break
		}

	}
	return new(T), errors.New("The allowed number of iterations has been exceeded.")
}
