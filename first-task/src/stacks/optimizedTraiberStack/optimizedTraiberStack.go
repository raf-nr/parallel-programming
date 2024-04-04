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
	return &Stack[T]{exchangerArray: freshExchangersArray[T](10, 5)}
}

func (stack *Stack[T]) IsEmpty() bool {
	return stack.top.Load() == nil
}

func (stack *Stack[T]) Peek() (T, error) {
	if stack.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New("Stack is already empty.")
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
		if stack.primitePush(value) {
			break
		}
		_, err := stack.exchangerArray.visit(&value)
		if err == nil {
			return nil
		}
	}
	return nil
}

func (stack *Stack[T]) primitivePop() (T, error) {
	if stack.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New(stacks.EmptyStackError)
	}
	oldTop := stack.top.Load()
	newTop := oldTop.next.Load()
	if stack.top.CompareAndSwap(oldTop, newTop) {
		return oldTop.value, nil
	}
	return *new(T), errors.New(stacks.UnsuccessfulPrimitivePop)
}

func (stack *Stack[T]) Pop() (T, error) {
	for {
		value, err := stack.primitivePop()
		if err == nil {
			return value, err
		}

		if err.Error() == stacks.UnsuccessfulPrimitivePop {
			element, err := stack.exchangerArray.visit(nil)
			if err == nil {
				return *element, nil
			}
		} else {
			return *new(T), err
		}
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
