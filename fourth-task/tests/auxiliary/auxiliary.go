package auxiliary

import (
	"bst/trees"
	"bst/trees/coarseGrainedTree"
	"bst/trees/fineGrainedTree"
	"bst/trees/optimisticTree"
)

func FreshCoarseGrainedTree() trees.BinarySearchTree[int, int] {
	return coarseGrainedTree.FreshCoarseGrainedSyncTree[int, int]()
}

func FreshFineGrainedTree() trees.BinarySearchTree[int, int] {
	return fineGrainedTree.FreshFineGrainedSyncTree[int, int]()
}

func FreshOptimisticTree() trees.BinarySearchTree[int, int] {
	return optimisticTree.FreshOptimisticSyncTree[int, int]()
}
