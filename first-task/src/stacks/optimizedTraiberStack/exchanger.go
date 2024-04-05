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

func (e *exchanger[T]) exchange(value *T, replays int) (*T, error) {

	for i := 0; i < replays; i++ {
		// "Knock" on the cell until its status becomes 'empty' or 'waiting'.
		currentElem := e.elem.Load().(exchangerElem[T])

		switch currentElem.state {
		case empty:
			// If the cell is empty, put a new value and wait until another thread arrives with a complementary operation.
			oldItem := currentElem
			newItem := exchangerElem[T]{value, waiting}
			if e.elem.CompareAndSwap(oldItem, newItem) {
				for j := i; j < replays; j++ {
					currentElem := e.elem.Load().(exchangerElem[T])
					if currentElem.state == busy {
						// The block is executed if a stream arrives with a complementary operation.
						newItem := exchangerElem[T]{nil, empty}
						e.elem.Store(newItem)
						return currentElem.value, nil // Make an exchange.
					}
				}
				e.elem.Store(oldItem)
				return new(T), errors.New("The allowed number of iterations has been exceeded.")
			}
		case waiting:
			// If the cell is occupied by a thread and is waiting for another,
			// check that the operations are complementary,
			// and if so, carry out the exchange
			oldItem := currentElem
			newItem := exchangerElem[T]{value, busy}
			if (oldItem.value == nil) != (value == nil) {
				if e.elem.CompareAndSwap(oldItem, newItem) {
					return currentElem.value, nil
				}
			}
		case busy:
			// In this case, an exchange has occurred in the cell at the moment,
			// so just wait until it becomes free
			break
		}

	}
	return new(T), errors.New("The allowed number of iterations has been exceeded.")
}
