package tests

import (
	"src/stacks"
	"src/tests/auxiliary"
	"testing"
)

// In these test cases, we run all types of tests sequentially.

const elementsAmount = 1_000_000

func TestConsistentStackSequential(t *testing.T) {
	runStackTests(t, auxiliary.FreshConsistentStack)
}

func TestTraiberStackSequential(t *testing.T) {
	runStackTests(t, auxiliary.FreshTraiberStack)
}

func TestOptimizedTraiberStackSequential(t *testing.T) {
	runStackTests(t, auxiliary.FreshOptimizedTraiberStack)
}

func runStackTests(t *testing.T, newStack func() stacks.Stack[int]) {

	t.Run("Test Empty Stack Pop: ", func(t *testing.T) {
		stack := newStack()
		elem, err := stack.Pop()

		if elem != 0 && err == nil {
			t.Errorf("Error: some value was received instead of the expected emptyStackError.")
		} else {
			if err.Error() != stacks.EmptyStackError {
				t.Errorf("Error: received an error other than the expected emptyStackError.")
			}
		}
	})

	t.Run("Test Empty Stack Peek: ", func(t *testing.T) {
		stack := newStack()
		elem, err := stack.Peek()
		if elem != 0 && err == nil {
			t.Errorf("Error: some value was received instead of the expected emptyStackError.")
		} else {
			if err.Error() != stacks.EmptyStackError {
				t.Errorf("Error: received an error other than the expected emptyStackError.")
			}
		}
	})

	t.Run("Test Not Empty Stack Peek: ", func(t *testing.T) {
		stack := newStack()
		stack.Push(1)
		stack.Push(2)
		elem, err := stack.Peek()
		expected := 2

		if elem != expected && err == nil {
			t.Errorf("Received top %d != expected top %d", elem, expected)
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("Test Stack Push and Pop: ", func(t *testing.T) {
		stack := newStack()
		for i := 0; i < 10; i++ {
			stack.Push(i)
		}

		for i := 0; i < 10; i++ {
			stack.Pop()
		}

		elem, err := stack.Pop()

		if elem != 0 && err == nil {
			t.Errorf("Error: some value was received instead of the expected emptyStackError.")
		} else {
			if err.Error() != stacks.EmptyStackError {
				t.Errorf("Error: received an error other than the expected emptyStackError.")
			}
		}
	})

	t.Run("Test Stack Push and Pop2: ", func(t *testing.T) {
		stack := newStack()
		stack.Push(1)
		stack.Push(2)
		stack.Push(3)
		stack.Pop()

		elem, err := stack.Pop()

		expected := 2

		if elem != expected && err == nil {
			t.Errorf("Received removed element %d != expected removed element %d", elem, expected)
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("Test Stack Len: ", func(t *testing.T) {
		stack := newStack()
		for i := 1; i <= elementsAmount; i++ {
			stack.Push(i)
		}

		stackLen, err := stack.Len()

		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if stackLen != elementsAmount {
			t.Errorf("Reveived stack len %d != expected stack len %d", stackLen, elementsAmount)
		}
	})
}
