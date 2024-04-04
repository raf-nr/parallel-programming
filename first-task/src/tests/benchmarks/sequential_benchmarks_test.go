package benchmarks

import (
	"math/rand"
	"src/stacks"
	"src/tests/auxiliary"
	"testing"
)

func BenchmarkSequentialConsistentStack(b *testing.B) {
	runsStackBenchmarks(b, auxiliary.FreshConsistentStack)
}

func BenchmarkSequentialTraiberStack(b *testing.B) {
	runsStackBenchmarks(b, auxiliary.FreshTraiberStack)
}

func BenchmarkSequentialOptimizedTraiberStack(b *testing.B) {
	runsStackBenchmarks(b, auxiliary.FreshOptimizedTraiberStack)
}

const elementsAmount = 1_000_000

func runsStackBenchmarks(b *testing.B, newStack func() stacks.Stack[int]) {

	b.Run("Push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			for j := 0; j < elementsAmount; j++ {
				stack.Push(j)
			}
		}
	})

	b.Run("Pop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			for j := 0; j < elementsAmount; j++ {
				stack.Pop()
			}
		}
	})

	b.Run("Push and pop in sequential order", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			for j := 0; j < elementsAmount; j++ {
				stack.Push(j)
				stack.Pop()
			}
		}
	})

	b.Run("Push and Pop in random order", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			for j := 0; j < elementsAmount; j++ {
				operation := rand.Intn(2)
				if operation == 0 {
					stack.Push(j)
				} else {
					stack.Pop()
				}
			}
		}
	})
}
