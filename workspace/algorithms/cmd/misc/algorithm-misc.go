package main

import (
	"fmt"
	"slices"
	"sync"
)

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

func main() {
	// ProblemIsSameTree1()
	// ProblemIsSameTree2()
	ProblemIsSameTree3()

}

// @problem: https://leetcode.com/problems/same-tree/
func ProblemIsSameTree1() {
	p := &Node{Val: 1, Left: &Node{Val: 2}, Right: &Node{Val: 3}}
	q := &Node{Val: 1, Left: &Node{Val: 2}, Right: &Node{Val: 3}}

	const NilMarker = 1<<31 - 1 // MaxInt, unlikely in input | we use this to mark nil nodes

	isSameTree := func(p *Node, q *Node) bool {
		var algorithm func(n *Node, traversal *[]int)

		algorithm = func(n *Node, traversal *[]int) {
			if n == nil {
				*traversal = append(*traversal, NilMarker)
				return
			}

			*traversal = append(*traversal, n.Val)
			algorithm(n.Left, traversal)
			algorithm(n.Right, traversal)
		}

		traversalP := make([]int, 0)
		traversalQ := make([]int, 0)

		algorithm(p, &traversalP)
		algorithm(q, &traversalQ)

		equal := slices.Equal(traversalP, traversalQ)
		fmt.Println(traversalP, traversalQ, equal)
		return equal
	}

	fmt.Println(
		`
==================================================
	ProblemIsSameTree1
==================================================`,
	)
	result := isSameTree(p, q)
	fmt.Println(result)

}

// @problem: https://leetcode.com/problems/same-tree/
func ProblemIsSameTree2() {
	p := &Node{Val: 1, Left: &Node{Val: 2}, Right: &Node{Val: 3}}
	q := &Node{Val: 1, Left: &Node{Val: 2}, Right: &Node{Val: 3}}

	const NilMarker = 1<<31 - 1 // MaxInt, unlikely in input | we use this to mark nil nodes

	isSameTree := func(p *Node, q *Node) bool {
		var wg sync.WaitGroup

		algorithm := func(n *Node, traversal *[]int) {
			var f func(n *Node, traversal *[]int)

			f = func(n *Node, traversal *[]int) {
				if n == nil {
					*traversal = append(*traversal, NilMarker)
					return
				}

				*traversal = append(*traversal, n.Val)
				f(n.Left, traversal)
				f(n.Right, traversal)
			}

			f(n, traversal)
			wg.Done()

		}

		traversalP := make([]int, 0)
		traversalQ := make([]int, 0)

		wg.Add(2)
		go algorithm(p, &traversalP)
		go algorithm(q, &traversalQ)

		wg.Wait()

		equal := slices.Equal(traversalP, traversalQ)
		fmt.Println(traversalP, traversalQ, equal)
		return equal
	}

	fmt.Println(
		`
==================================================
	ProblemIsSameTree2
==================================================`,
	)
	result := isSameTree(p, q)
	fmt.Println(result)

}

// @problem: https://leetcode.com/problems/same-tree/
func ProblemIsSameTree3() {
	p := &Node{Val: 1, Left: &Node{Val: 2}, Right: &Node{Val: 3}}
	q := &Node{Val: 1, Left: &Node{Val: 2}, Right: &Node{Val: 3}}

	var isSameTree func(p *Node, q *Node) bool
	isSameTree = func(p *Node, q *Node) bool {
		if p == nil || q == nil {
			return p == q
		}

		if p.Val != q.Val {
			return false
		}

		return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
	}

	fmt.Println(
		`
==================================================
	ProblemIsSameTree2
==================================================`,
	)
	result := isSameTree(p, q)
	fmt.Println(result)

}
