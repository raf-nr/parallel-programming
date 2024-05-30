package trees

import "cmp"

type BinarySearchTree[T any, K cmp.Ordered] interface {
	Find(K) (T, bool)
	Insert(K, T)
	Remove(K) bool
	CountNodes() int
	IsValid() bool
	Print()
}
