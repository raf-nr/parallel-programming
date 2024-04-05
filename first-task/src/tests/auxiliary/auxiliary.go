package auxiliary

import (
	"src/stacks"
	"src/stacks/consistentStack"
	"src/stacks/optimizedTraiberStack"
	"src/stacks/traiberStack"
)

// Helper functions that resolve type problems in a tests and benchmarks.

func FreshConsistentStack() stacks.Stack[int] {
	return consistentStack.FreshConsistentStack[int]()
}

func FreshTraiberStack() stacks.Stack[int] {
	return traiberStack.FreshTraiberStack[int]()
}

func FreshOptimizedTraiberStack() stacks.Stack[int] {
	return optimizedTraiberStack.FreshOptimizedTraiberStack[int]()
}
