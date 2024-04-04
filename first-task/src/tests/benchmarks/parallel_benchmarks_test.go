package benchmarks

import (
	"src/stacks"
	"src/tests/auxiliary"
	"sync"
	"testing"
)

func BenchmarkParallelTraiberStack(b *testing.B) {
	runsStackBenchmarksParallel(b, auxiliary.FreshTraiberStack)
}

func BenchmarkParallelOptimizedTraiberStack(b *testing.B) {
	runsStackBenchmarksParallel(b, auxiliary.FreshOptimizedTraiberStack)
}

func runsStackBenchmarksParallel(b *testing.B, newStack func() stacks.Stack[int]) {

	b.Run("Push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(elementsAmount)
			for j := 0; j < elementsAmount; j++ {
				go func() {
					defer wg.Done()
					stack.Push(j)
				}()
			}
			wg.Wait()
		}
	})

	b.Run("Pop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(elementsAmount)
			for j := 0; j < elementsAmount; j++ {
				go func() {
					defer wg.Done()
					stack.Pop()
				}()
			}
			wg.Wait()
		}
	})

	//b.Run("Push and pop in sequential order", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		stack := newStack()
	//		wg := sync.WaitGroup{}
	//		wg.Add(elementsAmount)
	//		for j := 0; j < elementsAmount; j++ {
	//			go func() {
	//				defer wg.Done()
	//				stack.Push(j)
	//				stack.Pop()
	//			}()
	//		}
	//		wg.Wait()
	//	}
	//})

	//b.Run("Push and pop in sequential order in different gorutines", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		stack := newStack()
	//		wg := sync.WaitGroup{}
	//		wg.Add(elementsAmount * 2)
	//		for j := 0; j < elementsAmount; j++ {
	//			go func() {
	//				defer wg.Done()
	//				stack.Push(j)
	//			}()
	//
	//			go func() {
	//				defer wg.Done()
	//				stack.Pop()
	//			}()
	//		}
	//		wg.Wait()
	//	}
	//})

	//b.Run("Push and Pop in random order", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		stack := newStack()
	//		wg := sync.WaitGroup{}
	//		wg.Add(elementsAmount)
	//		for j := 0; j < elementsAmount; j++ {
	//			operation := rand.Intn(2)
	//			if operation == 0 {
	//				go func() {
	//					defer wg.Done()
	//					stack.Push(j)
	//				}()
	//			} else {
	//				go func() {
	//					defer wg.Done()
	//					stack.Pop()
	//				}()
	//			}
	//		}
	//		wg.Wait()
	//	}
	//})
}
