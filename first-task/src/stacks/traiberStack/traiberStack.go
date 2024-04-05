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
	// New stack instance.
	return &Stack[T]{}
}

func (stack *Stack[T]) Peek() (T, error) {

	if stack == nil {
		return *(new(T)), errors.New(stacks.StackNilPointerError)
	}

	if stack.top.Load() == nil {
		return *(new(T)), errors.New(stacks.EmptyStackError)
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
	for {

		oldTop := stack.top.Load()
		if oldTop == nil {
			return *(new(T)), errors.New(stacks.EmptyStackError)
		}
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
