package consistentStack

import (
	"errors"
	"src/stacks"
)

type cell[T any] struct {
	value T
	next  *cell[T]
}

type Stack[T any] struct {
	top *cell[T]
}

func FreshConsistentStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (stack *Stack[T]) IsEmpty() bool {
	return stack.top == nil
}

func (stack *Stack[T]) Peek() (T, error) {
	if stack.top == nil {
		var zeroValue T
		return zeroValue, errors.New(stacks.EmptyStackError)
	}
	return stack.top.value, nil
}

func (stack *Stack[T]) Push(value T) error {
	if stack == nil {
		return errors.New(stacks.StackNilPointerError)
	}
	next := &cell[T]{value: value, next: stack.top}
	stack.top = next
	return nil
}

func (stack *Stack[T]) Pop() (T, error) {
	if stack == nil {
		var zeroValue T
		return zeroValue, errors.New(stacks.StackNilPointerError)
	}
	if stack.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New(stacks.EmptyStackError)
	}
	value := stack.top.value
	stack.top = stack.top.next
	return value, nil
}

func (stack *Stack[T]) Len() (int, error) {
	if stack == nil {
		return 0, errors.New(stacks.StackNilPointerError)
	}
	size := 0
	current := stack.top
	for current != nil {
		size++
		current = current.next
	}
	return size, nil
}
