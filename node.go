package concurrent

import "time"

type node struct {
	key        Comparable
	value      string
	left       *node
	right      *node
	leftDepth  int
	rightDepth int
	updateTime time.Time
}

func (t *node) Clone() *node {
	return &node{
		key:        t.key,
		value:      t.value,
		left:       t.left,
		right:      t.right,
		leftDepth:  t.leftDepth,
		rightDepth: t.rightDepth,
		updateTime: time.Now(),
	}
}

func createNode(key Comparable, value string) *node {
	return &node{
		key:        key,
		value:      value,
		left:       nil,
		right:      nil,
		updateTime: time.Now(),
	}
}

func updateDepth(n *node) {
	if n == nil {
		return
	}
	if n.left != nil {
		n.leftDepth = max(n.left.leftDepth, n.left.rightDepth) + 1
	} else {
		n.leftDepth = 0
	}
	if n.right != nil {
		n.rightDepth = max(n.right.rightDepth, n.right.rightDepth) + 1
	} else {
		n.rightDepth = 0
	}
}

func rotate(n *node) *node {
	if n.leftDepth > n.rightDepth && n.leftDepth-n.rightDepth >= 2 {
		return rightRotate(n)
	}
	if n.rightDepth > n.leftDepth && n.rightDepth-n.leftDepth >= 2 {
		return leftRotate(n)
	}
	return n
}

func leftRotate(n *node) *node {
	newRoot := n.right.Clone()
	newLeft := n.Clone()
	newLeft.right = newRoot.left
	newRoot.left = newLeft
	updateDepth(newRoot.left)
	updateDepth(newRoot.right)
	updateDepth(newRoot)
	return newRoot
}

func rightRotate(n *node) *node {
	newRoot := n.left.Clone()
	newRight := n.Clone()
	newRight.left = newRoot.right
	newRoot.right = newRight
	updateDepth(newRoot.left)
	updateDepth(newRoot.right)
	updateDepth(newRoot)
	return newRoot
}
