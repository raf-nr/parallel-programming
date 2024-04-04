package tests

import (
	"src/stacks"
	"src/tests/auxiliary"
	"testing"
)

const elements = 1_000_000

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

		if elem != 2 && err == nil {
			t.Errorf("Received top %d != expected top 2", elem)
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
			if elem != 0 && err == nil {
				t.Errorf("Error: some value was received instead of the expected emptyStackError.")
			} else {
				if err.Error() != stacks.EmptyStackError {
					t.Errorf("Error: received an error other than the expected emptyStackError.")
				}
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

		if elem != 2 && err == nil {
			t.Errorf("Received removed element %d != expected removed element 2", elem)
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("Test Stack Len: ", func(t *testing.T) {
		stack := newStack()
		for i := 1; i <= elements; i++ {
			stack.Push(i)
		}

		stackLen, err := stack.Len()

		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if stackLen != elements {
			t.Errorf("Reveived stack len %d != expected stack len %d", stackLen, elements)
		}
	})
}
