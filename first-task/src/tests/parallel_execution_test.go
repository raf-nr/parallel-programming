package tests

import (
	"src/stacks"
	"src/tests/auxiliary"
	"sync"
	"testing"
)

func TestTraiberStackParalell(t *testing.T) {
	runParallelStackTests(t, auxiliary.FreshTraiberStack)
}

func TestOptimizedTraiberStackParallel(t *testing.T) {
	runParallelStackTests(t, auxiliary.FreshOptimizedTraiberStack)
}

func runParallelStackTests(t *testing.T, newStack func() stacks.Stack[int]) {

	elementsAmount := 1_000_000

	t.Run("Test push", func(t *testing.T) {
		stack := newStack()
		wg := sync.WaitGroup{}
		wg.Add(elementsAmount)
		for i := 0; i < elementsAmount; i++ {
			go func() {
				defer wg.Done()
				err := stack.Push(i)
				if err != nil {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			}()
		}
		wg.Wait()

		stackLen, err := stack.Len()
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if stackLen != 1_000_000 {
			t.Errorf("Received stack size %d != expected stack size 1000000", stackLen)
		}
	})

	t.Run("Test pop on empty stack", func(t *testing.T) {
		stack := newStack()
		wg := sync.WaitGroup{}
		wg.Add(elementsAmount)
		for i := 0; i < elementsAmount; i++ {
			go func() {
				defer wg.Done()
				_, err := stack.Pop()
				if err.Error() != stacks.EmptyStackError {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			}()
		}
		wg.Wait()

		stackLen, err := stack.Len()
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if stackLen != 0 {
			t.Errorf("Received stack size %d != expected stack size 0", stackLen)
		}
	})

	t.Run("Test push and pop", func(t *testing.T) {
		stack := newStack()
		wg := sync.WaitGroup{}
		wg.Add(elementsAmount)
		for i := 0; i < elementsAmount; i++ {
			go func() {
				defer wg.Done()
				err := stack.Push(i)
				if err != nil {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			}()
		}
		wg.Wait()

		wg.Add(elementsAmount)
		for i := 0; i < elementsAmount; i++ {
			go func() {
				defer wg.Done()
				_, err := stack.Pop()
				if err != nil {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			}()
		}
		wg.Wait()
		stackLen, err := stack.Len()
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if stackLen != 0 {
			t.Errorf("Received stack size %d != expected stack size 0", stackLen)
		}
	})

}
