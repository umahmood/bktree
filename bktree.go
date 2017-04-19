package main

import (
	"container/list"
	"fmt"

	"github.com/arbovm/levenshtein"
)

// DistanceFunc a distance function which calculates the distance between two
// strings. This distance function must satisfy a set of axioms in order to
// ensure it's well-behaved. See https://en.wikipedia.org/wiki/Metric_space.
type DistanceFunc func(a, b string) int

// Results list of results
type Results []Result

// Result holds a single result
type Result struct {
	Distance int
	Value    string
}

// Node represents a single node in the BK-Tree
type Node struct {
	Value    string
	Children map[int]*Node
}

// BKTree a data structure specialized to index data in a metric space
type BKTree struct {
	Root     *Node
	distFunc DistanceFunc
}

// Add adds a value into the tree
func (t *BKTree) Add(node string) {
	if t.Root == nil {
		t.Root = &Node{Value: node, Children: make(map[int]*Node)}
		return
	}
	current, children := t.Root.Value, t.Root.Children
	for {
		dist := t.distFunc(node, current)
		target := children[dist]
		if target == nil {
			children[dist] = &Node{Value: node, Children: make(map[int]*Node)}
			break
		}
		current, children = target.Value, target.Children
	}
}

// Search the tree and return all words closest to a given query word.
func (t *BKTree) Search(node string, radius int) Results {
	if t.Root == nil {
		return nil
	}
	var results Results
	candidates := list.New()
	candidates.PushBack(t.Root)
	for e := candidates.Front(); e != nil; e = e.Next() {
		v := e.Value.(*Node)
		candidate, children := v.Value, v.Children
		dist := t.distFunc(node, candidate)
		if dist <= radius {
			results = append(results, Result{Distance: dist, Value: candidate})
		}
		low, high := dist-radius, dist+radius
		for d, c := range children {
			if low <= d && d <= high {
				candidates.PushBack(c)
			}
		}
	}
	return results
}

// NewBKTree creates a new BK-Tree instance
func NewBKTree(distFunc DistanceFunc) *BKTree {
	return &BKTree{distFunc: distFunc}
}

func main() {
	tree := NewBKTree(levenshtein.Distance)

	words := []string{"some", "soft", "same", "mole", "soda", "salmon"}

	for _, w := range words {
		tree.Add(w)
	}

	for _, result := range tree.Search("bole", 2) {
		fmt.Println("Value:", result.Value, "Distance:", result.Distance)
	}
}
