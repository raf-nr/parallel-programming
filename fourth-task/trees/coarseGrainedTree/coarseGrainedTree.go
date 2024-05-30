package coarseGrainedTree

import (
	"cmp"
	"fmt"
	"strings"
	"sync"
)

type CoarseGrainedSyncTree[T any, K cmp.Ordered] struct {
	root  *Node[T, K]
	mutex sync.Mutex
}

type Node[T any, K cmp.Ordered] struct {
	key   K
	value T
	left  *Node[T, K]
	right *Node[T, K]
}

func FreshCoarseGrainedSyncTree[T any, K cmp.Ordered]() *CoarseGrainedSyncTree[T, K] {
	return &CoarseGrainedSyncTree[T, K]{
		root:  nil,
		mutex: sync.Mutex{},
	}
}

func (tree *CoarseGrainedSyncTree[T, K]) lock() {
	tree.mutex.Lock()
}

func (tree *CoarseGrainedSyncTree[T, K]) unlock() {
	tree.mutex.Unlock()
}

func (tree *CoarseGrainedSyncTree[T, K]) Insert(key K, value T) {
	tree.lock()
	defer tree.unlock()
	if tree.root == nil {
		tree.root = &Node[T, K]{key: key, value: value}
		return
	}
	tree.insertRecursive(tree.root, key, value)
}

func (tree *CoarseGrainedSyncTree[T, K]) insertRecursive(node *Node[T, K], key K, value T) {
	if cmp.Less(key, node.key) {
		if node.left == nil {
			node.left = &Node[T, K]{key: key, value: value}
		} else {
			tree.insertRecursive(node.left, key, value)
		}
	} else if cmp.Compare(key, node.key) == 1 {
		if node.right == nil {
			node.right = &Node[T, K]{key: key, value: value}
		} else {
			tree.insertRecursive(node.right, key, value)
		}
	} else {
		// Key already exists, update value
		node.value = value
	}
}

func (tree *CoarseGrainedSyncTree[T, K]) Find(key K) (T, bool) {
	tree.lock()
	defer tree.unlock()
	value, found := tree.findRecursive(tree.root, key)
	return value, found
}

func (tree *CoarseGrainedSyncTree[T, K]) findRecursive(node *Node[T, K], key K) (T, bool) {
	var emptValue T
	if node == nil {
		return emptValue, false
	}
	comp := cmp.Compare(key, node.key)
	if comp == -1 {
		return tree.findRecursive(node.left, key)
	} else if comp == 1 {
		return tree.findRecursive(node.right, key)
	} else {
		return node.value, true
	}
}

func (tree *CoarseGrainedSyncTree[T, K]) Remove(key K) bool {
	tree.lock()
	defer tree.unlock()
	_, removed := tree.removeRecursive(tree.root, key)
	return removed
}

func (tree *CoarseGrainedSyncTree[T, K]) removeRecursive(node *Node[T, K], key K) (*Node[T, K], bool) {
	if node == nil {
		return nil, false
	}
	var removed bool
	comp := cmp.Compare(key, node.key)
	if comp == -1 {
		node.left, removed = tree.removeRecursive(node.left, key)
		return node, removed
	} else if comp == 1 {
		node.right, removed = tree.removeRecursive(node.right, key)
		return node, removed
	} else {
		// Node to be deleted found
		if node.left == nil {
			return node.right, true
		} else if node.right == nil {
			return node.left, true
		} else {
			// Node to be deleted has two children
			// Find the inorder successor (smallest in the right subtree)
			successor := tree.minValueNode(node.right)
			// Copy the inorder successor's content to this node
			node.key = successor.key
			node.value = successor.value
			// Delete the inorder successor
			node.right, _ = tree.removeRecursive(node.right, successor.key)
			return node, true
		}
	}
}

func (tree *CoarseGrainedSyncTree[T, K]) minValueNode(node *Node[T, K]) *Node[T, K] {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

func (tree *CoarseGrainedSyncTree[T, K]) Print() {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()

	tree.printHelper(tree.root, 0)
}

func (tree *CoarseGrainedSyncTree[T, K]) printHelper(node *Node[T, K], indent int) {
	if node == nil {
		return
	}

	// Print the right subtree with indentation
	tree.printHelper(node.right, indent+4)

	// Print current node
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Printf("%v\n", node.key)

	// Print left subtree with indentation
	tree.printHelper(node.left, indent+4)
}

func (tree *CoarseGrainedSyncTree[T, K]) CountNodes() int {
	// CountNodes counts the number of nodes in the tree
	tree.lock()
	defer tree.unlock()
	return tree.countNodesRecursive(tree.root)
}

func (tree *CoarseGrainedSyncTree[T, K]) countNodesRecursive(node *Node[T, K]) int {
	if node == nil {
		return 0
	}
	return 1 + tree.countNodesRecursive(node.left) + tree.countNodesRecursive(node.right)
}

func (tree *CoarseGrainedSyncTree[T, K]) IsValid() bool {
	// IsValid checks if the tree is a valid binary search tree
	tree.lock()
	defer tree.unlock()
	return tree.isValidBSTRecursive(tree.root, nil, nil)
}

func (tree *CoarseGrainedSyncTree[T, K]) isValidBSTRecursive(node *Node[T, K], min *K, max *K) bool {
	if node == nil {
		return true
	}
	if (min != nil && cmp.Compare(node.key, *min) <= 0) || (max != nil && cmp.Compare(node.key, *max) >= 0) {
		return false
	}
	return tree.isValidBSTRecursive(node.left, min, &node.key) && tree.isValidBSTRecursive(node.right, &node.key, max)
}
