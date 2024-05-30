package optimisticTree

import (
	"cmp"
	"fmt"
	"strings"
	"sync"
)

type OptimisticSyncTree[T any, K cmp.Ordered] struct {
	root  *Node[T, K]
	mutex sync.Mutex
}

type Node[T any, K cmp.Ordered] struct {
	key   K
	value T
	left  *Node[T, K]
	right *Node[T, K]
	mutex sync.Mutex
}

func FreshOptimisticSyncTree[T any, K cmp.Ordered]() *OptimisticSyncTree[T, K] {
	return &OptimisticSyncTree[T, K]{
		root:  nil,
		mutex: sync.Mutex{},
	}
}

func (tree *OptimisticSyncTree[T, K]) lock() {
	tree.mutex.Lock()
}

func (tree *OptimisticSyncTree[T, K]) unlock() {
	tree.mutex.Unlock()
}

func (node *Node[T, K]) lock() {
	node.mutex.Lock()
}

func (node *Node[T, K]) unlock() {
	node.mutex.Unlock()
}

func (tree *OptimisticSyncTree[T, K]) findWithParent(key K) (*Node[T, K], *Node[T, K]) {
	for {

		tree.lock()
		if tree.root == nil {
			return nil, nil
		}

		parent := (*Node[T, K])(nil)
		current := tree.root

		// Try to find a node without locks
		for current != nil && current.key != key {
			if parent == nil {
				tree.unlock()
			}
			parent = current
			if cmp.Less(key, current.key) {
				current = current.left
			} else {
				current = current.right
			}
		}

		// Locating the found nodes
		if parent != nil {
			parent.lock()
		}

		if current != nil {
			current.lock()
		}

		// Validate the path
		validateParent := (*Node[T, K])(nil)
		validateCurrent := tree.root

		for validateCurrent != nil && validateCurrent != current && cmp.Compare(validateCurrent.key, key) != 0 {
			validateParent = validateCurrent
			if cmp.Less(key, validateCurrent.key) {
				validateCurrent = validateCurrent.left
			} else {
				validateCurrent = validateCurrent.right
			}
		}

		// Checking the validity of current and parent nodes
		if validateCurrent != current || validateParent != parent {
			if current != nil {
				current.unlock()
			}
			if parent != nil {
				parent.unlock()
			}
			continue
		}

		return current, parent
	}
}

func (tree *OptimisticSyncTree[T, K]) Find(key K) (T, bool) {
	// Use a helper method for optimistic searching
	node, parent := tree.findWithParent(key)

	// If the node is not found
	if node == nil {
		if parent != nil {
			parent.unlock()
		} else {
			tree.unlock()
		}
		var zeroValue T
		return zeroValue, false
	}

	// Read the value from the node
	value := node.value

	// Remove locks
	node.unlock()
	if parent != nil {
		parent.unlock()
	} else {
		tree.unlock()
	}

	return value, true
}

func (tree *OptimisticSyncTree[T, K]) Insert(key K, value T) {
	// Use the findWithParent helper method to find a node and its parent
	node, parent := tree.findWithParent(key)

	// If a node with such a key already exists, update its value
	if node != nil {
		node.value = value
		node.unlock()
		if parent != nil {
			parent.unlock()
		} else {
			tree.unlock()
		}
		return
	}

	newNode := &Node[T, K]{
		key:   key,
		value: value,
		left:  nil,
		right: nil,
	}

	// Insert a new node into the tree
	if parent == nil {
		// The tree is empty, insert a new node as the root
		tree.root = newNode
		tree.unlock()
	} else {
		// Insert a new node as a child of the parent
		if cmp.Less(key, parent.key) {
			parent.left = newNode
		} else {
			parent.right = newNode
		}
		parent.unlock()
	}

	return
}

func (tree *OptimisticSyncTree[T, K]) Remove(key K) bool {
	// Use the findWithParent helper method to find a node and its parent
	node, parent := tree.findWithParent(key)

	// If the node is not found, do nothing
	if node == nil {
		if parent != nil {
			parent.unlock()
		} else {
			tree.unlock()
		}
		return false
	}

	defer func() {
		node.unlock()
		if parent != nil {
			parent.unlock()
		} else {
			tree.unlock()
		}
	}()

	// If the node has no children
	if node.left == nil && node.right == nil {
		if parent == nil {
			tree.root = nil
		} else if cmp.Less(key, parent.key) {
			parent.left = nil
		} else {
			parent.right = nil
		}
		return true
	}

	// If the node has one child
	if node.left == nil || node.right == nil {
		var child *Node[T, K]
		if node.left != nil {
			child = node.left
		} else {
			child = node.right
		}

		if parent == nil {
			tree.root = child
		} else if cmp.Less(key, parent.key) {
			parent.left = child
		} else {
			parent.right = child
		}
		return true
	}

	// If the node has two children
	if node.left != nil && node.right != nil {
		successorParent := node
		successor := node.right
		successor.lock()

		for successor.left != nil {
			successor.left.lock()
			if successorParent != node {
				successorParent.unlock()
			}
			successorParent = successor
			successor = successor.left
		}

		if successorParent != node {
			successorParent.left = successor.right
			node.key = successor.key
			node.value = successor.value
			successorParent.unlock()
		} else {
			successorParent.right = successor.right
			node.key = successor.key
			node.value = successor.value
		}

		successor.unlock()
		return true
	}

	return false
}

func (tree *OptimisticSyncTree[T, K]) Print() {
	tree.lock()
	defer tree.unlock()

	tree.printHelper(tree.root, 0)
}

func (tree *OptimisticSyncTree[T, K]) printHelper(node *Node[T, K], indent int) {
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

func (tree *OptimisticSyncTree[T, K]) CountNodes() int {
	// CountNodes counts the number of nodes in the tree
	tree.lock()
	defer tree.unlock()
	return tree.countNodesRecursive(tree.root)
}

func (tree *OptimisticSyncTree[T, K]) countNodesRecursive(node *Node[T, K]) int {
	if node == nil {
		return 0
	}
	return 1 + tree.countNodesRecursive(node.left) + tree.countNodesRecursive(node.right)
}

func (tree *OptimisticSyncTree[T, K]) IsValid() bool {
	// IsValid checks if the tree is a valid binary search tree
	tree.lock()
	defer tree.unlock()
	return tree.isValidBSTRecursive(tree.root, nil, nil)
}

func (tree *OptimisticSyncTree[T, K]) isValidBSTRecursive(node *Node[T, K], min *K, max *K) bool {
	if node == nil {
		return true
	}
	if (min != nil && cmp.Compare(node.key, *min) <= 0) || (max != nil && cmp.Compare(node.key, *max) >= 0) {
		return false
	}
	return tree.isValidBSTRecursive(node.left, min, &node.key) && tree.isValidBSTRecursive(node.right, &node.key, max)
}
