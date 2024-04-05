package benchmarks

import (
	"math/rand"
	"runtime"
	"src/stacks"
	"src/tests/auxiliary"
	"sync"
	"testing"
)

// Metrics are measured for parallel operations with thread-safe stacks.

func BenchmarkParallelTraiberStack(b *testing.B) {
	runtime.GOMAXPROCS(16)
	runParallelBenchmarks(b, auxiliary.FreshTraiberStack)
}

func BenchmarkParallelOptimizedTraiberStack(b *testing.B) {
	runtime.GOMAXPROCS(16)
	runParallelBenchmarks(b, auxiliary.FreshOptimizedTraiberStack)
}

const gorutinesAmount1 = 8
const gorutinesAmount2 = 100

func runParallelBenchmarks(b *testing.B, newStack func() stacks.Stack[int]) {

	b.Run("Push | All gorutines", func(b *testing.B) {
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

	b.Run("Pop | All gorutines", func(b *testing.B) {
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

	b.Run("Push and pop in sequential order | All gorutines", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(elementsAmount)
			for j := 0; j < elementsAmount; j++ {
				go func() {
					defer wg.Done()
					stack.Push(j)
					stack.Pop()
				}()
			}
			wg.Wait()
		}
	})

	b.Run("Push and pop in sequential order | 8 gorutines", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(gorutinesAmount1)
			for j := 0; j < gorutinesAmount1; j++ {
				go func() {
					defer wg.Done()
					for j := 0; j < elementsAmount/gorutinesAmount1; j++ {
						stack.Push(j)
						stack.Pop()
					}
				}()
			}
			wg.Wait()
		}
	})

	b.Run("Push and pop in sequential order | 100 gorutines", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(gorutinesAmount2)
			for j := 0; j < gorutinesAmount2; j++ {
				go func() {
					defer wg.Done()
					for j := 0; j < elementsAmount/gorutinesAmount2; j++ {
						stack.Push(j)
						stack.Pop()
					}
				}()
			}
			wg.Wait()
		}
	})

	b.Run("Push and pop in sequential order in different gorutines | All gorutines", func(b *testing.B) {
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

	b.Run("Push and Pop in random order | All gorutines", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(elementsAmount)
			for j := 0; j < elementsAmount; j++ {
				operation := rand.Intn(2)
				if operation == 0 {
					go func() {
						defer wg.Done()
						stack.Push(j)
					}()
				} else {
					go func() {
						defer wg.Done()
						stack.Pop()
					}()
				}
			}
			wg.Wait()
		}
	})

	b.Run("Push and Pop in random order | 8 gorutines", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(gorutinesAmount1)
			for j := 0; j < gorutinesAmount1; j++ {
				if rand.Intn(2) == 0 {
					go func() {
						defer wg.Done()
						for j := 0; j < elementsAmount/gorutinesAmount1; j++ {
							stack.Push(j)
						}
					}()
				} else {
					go func() {
						defer wg.Done()
						for j := 0; j < elementsAmount/gorutinesAmount1; j++ {
							stack.Pop()
						}
					}()
				}
			}
			wg.Wait()
		}
	})

	b.Run("Push and Pop in random order | 100 gorutines", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack := newStack()
			wg := sync.WaitGroup{}
			wg.Add(gorutinesAmount2)
			for j := 0; j < gorutinesAmount2; j++ {
				if rand.Intn(2) == 0 {
					go func() {
						defer wg.Done()
						for j := 0; j < elementsAmount/gorutinesAmount2; j++ {
							stack.Push(j)
						}
					}()
				} else {
					go func() {
						defer wg.Done()
						for j := 0; j < elementsAmount/gorutinesAmount2; j++ {
							stack.Pop()
						}
					}()
				}
			}
			wg.Wait()
		}
	})
}
