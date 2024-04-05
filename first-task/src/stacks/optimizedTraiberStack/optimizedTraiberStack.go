package optimizedTraiberStack

import (
	"errors"
	"src/stacks"
	"sync/atomic"
)

type cell[T any] struct {
	value T
	next  atomic.Pointer[cell[T]]
}

type Stack[T any] struct {
	top            atomic.Pointer[cell[T]]
	exchangerArray exchangersArray[T]
}

func FreshOptimizedTraiberStack[T any]() *Stack[T] {
	// New stack instance.
	return &Stack[T]{exchangerArray: freshExchangersArray[T](10, 500)}
}

func (stack *Stack[T]) Peek() (T, error) {

	if stack == nil {
		return *(new(T)), errors.New("The consistentStack pointer is nil.")
	}

	if stack.top.Load() == nil {
		return *(new(T)), errors.New("Stack is already empty.")
	}

	return stack.top.Load().value, nil
}

func (stack *Stack[T]) primitePush(value T) bool {
	newTop := &cell[T]{value: value}
	oldTop := stack.top.Load()
	newTop.next.Store(oldTop)
	return stack.top.CompareAndSwap(oldTop, newTop)
}

func (stack *Stack[T]) Push(value T) error {
	if stack == nil {
		return errors.New("The consistentStack pointer is nil.")
	}
	for {
		if stack.primitePush(value) { // Try to push the element.
			break
		}
		// If it was not possible to push an element,
		// put it in the array of exchangers and try to carry out the exchange.
		_, err := stack.exchangerArray.visit(&value)
		if err == nil {
			return nil
		}
		// If the exchange also fails - start over.
	}
	return nil
}

func (stack *Stack[T]) primitivePop() (T, error) {
	oldTop := stack.top.Load()
	if oldTop == nil {
		var zeroValue T
		return zeroValue, errors.New(stacks.EmptyStackError)
	}
	newTop := oldTop.next.Load()
	if stack.top.CompareAndSwap(oldTop, newTop) {
		return oldTop.value, nil
	}
	return *new(T), errors.New(stacks.UnsuccessfulPrimitivePop)
	// If it was not possible to delete,
	// will give an error that will tell that it is worth trying to make an exchange.
}

func (stack *Stack[T]) Pop() (T, error) {
	for {
		value, err := stack.primitivePop() // Try to remove the element.
		if err == nil {
			return value, err
		}

		if err.Error() == stacks.UnsuccessfulPrimitivePop {
			// If it was not possible to delete an element,
			// put it in the array of exchangers and try to carry out the exchange.
			element, err := stack.exchangerArray.visit(nil)
			if err == nil {
				return *element, nil
			}
		} else {
			return *new(T), err
		}
		// If the exchange also fails - start over.
	}
}

func (stack *Stack[T]) Len() (int, error) {
	if stack == nil {
		return 0, errors.New(stacks.StackNilPointerError)
	}
	size := 0
	current := stack.top.Load()
	for current != nil {
		size++
		current = current.next.Load()
	}
	return size, nil
}
