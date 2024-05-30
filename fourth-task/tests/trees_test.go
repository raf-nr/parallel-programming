package tests

import (
	"bst/tests/auxiliary"
	"bst/trees"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCoarseGrainedTree(t *testing.T) {
	runTreesTests(t, auxiliary.FreshCoarseGrainedTree)
}

func TestFineGrainedTree(t *testing.T) {
	runTreesTests(t, auxiliary.FreshFineGrainedTree)
}

func TestOptimisticTree(t *testing.T) {
	runTreesTests(t, auxiliary.FreshOptimisticTree)
}

func runTreesTests(t *testing.T, newTree func() trees.BinarySearchTree[int, int]) {

	t.Run("Test insert", func(t *testing.T) {
		/* This test runs 100 goroutines, each of which performs an insert.
		And then it checks that all values have been inserted and the tree
		structure is not broken. */
		const nodesAmount = 100

		tree := newTree()
		wg := sync.WaitGroup{}
		wg.Add(nodesAmount)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				tree.Insert(i, i*i)

			}(i)
		}
		wg.Wait()

		if !tree.IsValid() {
			t.Errorf("Error: tree is not valid.")
		}

		if tree.CountNodes() != nodesAmount {
			t.Errorf("Error: the tree contains %d nodes, although %d were expected", tree.CountNodes(), nodesAmount)
		}

	})

	t.Run("Test find after insert | First", func(t *testing.T) {
		/* This test starts 50 insert goroutines, waits for them to execute,
		and then runs 50 search goroutines, which check that all goroutines have been
		inserted before.
		Then the tree structure and the number of inserted elements are checked. */
		tree := newTree()

		_, flag := tree.Find(1)

		if flag {
			t.Errorf("The find function found a non-existent node in the tree")
		}
		const nodesAmount = 50
		wg := sync.WaitGroup{}
		wg.Add(nodesAmount)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				tree.Insert(i, i*i)

			}(i)
		}
		wg.Wait()

		wg.Add(nodesAmount)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				value, flag := tree.Find(i)

				if !flag {
					t.Errorf("The find function did not find an existing node.")
				}

				if value != i*i {
					t.Errorf("The node found contains the value %d when the value %d was expected.", value, i*i)
				}
			}(i)
		}
		wg.Wait()

		if !tree.IsValid() {
			t.Errorf("The find function broke the tree")
		}
		if tree.CountNodes() != nodesAmount {
			t.Errorf("The find function change the tree")
		}
	})

	t.Run("Test find after insert | Second", func(t *testing.T) {
		/* 300 goroutines are launched, each of which performs the insertion and
		immediately checks for the presence of the inserted element in the tree.
		The test checks the correctness of the find function itself. */
		tree := newTree()

		const nodesAmount = 300
		wg := sync.WaitGroup{}
		wg.Add(nodesAmount)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				tree.Insert(i, i*i)
				value, flag := tree.Find(i)

				if !flag {
					t.Errorf("The find function did not find an existing node.")
				}

				if value != i*i {
					t.Errorf("The node found contains the value %d when the value %d was expected.", value, i*i)
				}
			}(i)
		}
		wg.Wait()
		if !tree.IsValid() {
			t.Errorf("The find function broke the tree")
		}
		if tree.CountNodes() != nodesAmount {
			t.Errorf("The find function change the tree")
		}
	})

	t.Run("Test find and insert parallel", func(t *testing.T) {
		/* The test runs 600 goroutines - 300 for insertion and 300 for search.
		Then it checks that the tree was built correctly and that it actually has 300 nodes.
		The test essentially checks that operations do not conflict with each other and that there
		is no data race. */
		tree := newTree()

		const nodesAmount = 300
		wg := sync.WaitGroup{}
		wg.Add(nodesAmount * 2)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				tree.Insert(i, i*i)
			}(i)
			go func(i int) {
				defer wg.Done()
				value, flag := tree.Find(i)

				if flag && value != i*i {
					t.Errorf("The node found contains the value %d when the value %d was expected.", value, i*i)
				}
			}(i)
		}
		wg.Wait()
		if !tree.IsValid() {
			t.Errorf("The find function broke the tree")
		}
		if tree.CountNodes() != nodesAmount {
			t.Errorf("The find function change the tree")
		}
	})

	t.Run("Test remove | First", func(t *testing.T) {
		/* The test first inserts 500 elements into the tree,
		and then runs 500 goroutines to remove these elements.
		Then it checks that the tree is empty. Essentially,
		the test verifies that remove works correctly and is thread safe. */
		const nodesAmount = 500

		tree := newTree()

		wg := sync.WaitGroup{}
		wg.Add(nodesAmount)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				tree.Insert(i, i*i)
			}(i)
		}
		wg.Wait()

		wg.Add(nodesAmount)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				flag := tree.Remove(i)
				if !flag {
					t.Errorf("Failed to remove a node that was previously added.")
				}
			}(i)
		}
		wg.Wait()

		if !tree.IsValid() && tree.CountNodes() != 0 {
			t.Errorf("The tree was expected to be empty.")
		}
	})

	t.Run("Test remove | Second", func(t *testing.T) {
		/* The test runs 1000 goroutines - 500 for insertion and
		500 for deletion at the same time. Then it checks that the tree
		structure is not broken. The test essentially checks that operations do not
		conflict with each other and that data races do not occur. */
		const nodesAmount = 500

		tree := newTree()

		wg := sync.WaitGroup{}
		wg.Add(nodesAmount * 2)
		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				tree.Insert(i, i*i)
			}(i)
			go func(i int) {
				defer wg.Done()
				tree.Remove(i)
			}(i)
		}
		wg.Wait()
		if !tree.IsValid() {
			t.Errorf("Remove broke tree.")
		}
	})

	t.Run("Test insert and remove in random order", func(t *testing.T) {
		/* The test runs 500 goroutines, which with equal probability perform either
		insertion or deletion of an element. Essentially, the test checks that functions do not
		conflict with each other and that data races do not occur. */
		const nodesAmount = 500

		tree := newTree()

		wg := sync.WaitGroup{}
		wg.Add(nodesAmount)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 0; i < nodesAmount; i++ {
			go func(i int) {
				defer wg.Done()
				if r.Intn(2) == 0 {
					tree.Insert(i, i*i)
				} else {
					tree.Remove(i)
				}
			}(i)
		}

		wg.Wait()
		if !tree.IsValid() {
			t.Errorf("Insert and Remove broke tree.")
		}
	})

	t.Run("Test isValid", func(t *testing.T) {
		/* The test builds a valid tree and
		then checks that the isValid function works correctly. */
		tree := newTree()
		if !tree.IsValid() {
			t.Errorf("Error: function IsValid is not working correctly")
		}
		tree.Insert(10, 1)
		tree.Insert(-5, 2)
		tree.Insert(12, 1)
		tree.Insert(11, 3)
		tree.Insert(-7, 3)
		tree.Insert(-6, 3)
		tree.Insert(-4, 3)
		tree.Insert(-3, 3)
		if !tree.IsValid() {
			t.Errorf("Error: function IsValid is not working correctly")
		}
	})

}
