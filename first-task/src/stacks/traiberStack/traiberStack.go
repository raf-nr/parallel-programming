package traiberStack

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
	top atomic.Pointer[cell[T]]
}

func FreshTraiberStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (stack *Stack[T]) IsEmpty() bool {
	return stack.top.Load() == nil
}

func (stack *Stack[T]) Peek() (T, error) {
	if stack.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New(stacks.EmptyStackError)
	}
	return stack.top.Load().value, nil
}

func (stack *Stack[T]) Push(value T) error {
	if stack == nil {
		return errors.New(stacks.StackNilPointerError)
	}
	newTop := &cell[T]{value: value}
	for {
		oldTop := stack.top.Load()
		newTop.next.Store(oldTop)
		if stack.top.CompareAndSwap(oldTop, newTop) {
			return nil
		}
	}
}

func (stack *Stack[T]) Pop() (T, error) {
	if stack == nil {
		var zeroValue T
		return zeroValue, errors.New(stacks.StackNilPointerError)
	}
	for {
		if stack.IsEmpty() {
			var zeroValue T
			return zeroValue, errors.New(stacks.EmptyStackError)
		}
		oldTop := stack.top.Load()
		newTop := oldTop.next.Load()
		if stack.top.CompareAndSwap(oldTop, newTop) {
			return oldTop.value, nil
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
