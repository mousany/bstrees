package rb

import "golang.org/x/exp/constraints"

type rbColor bool

const (
	red   rbColor = true
	black rbColor = false
)

type rbTreeNode[T constraints.Ordered] struct {
	value T
	left  *rbTreeNode[T]
	right *rbTreeNode[T]
	color rbColor
	size  uint // Size of subtree, unnecessary if you don't need kth element
	// Father *RBNode[T] // Not necessary, but easier to implement
}

func newRBTreeNode[T constraints.Ordered](value T) *rbTreeNode[T] {
	return &rbTreeNode[T]{value: value, left: nil, right: nil, color: red, size: 1}
}

func (n *rbTreeNode[T]) Update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

func (n *rbTreeNode[T]) red() bool {
	return n.color == red
}

func (n *rbTreeNode[T]) child(direction bool) *rbTreeNode[T] {
	if direction {
		return n.right
	} else {
		return n.left
	}
}

func (n *rbTreeNode[T]) setChild(direction bool, child *rbTreeNode[T]) {
	if direction {
		n.right = child
	} else {
		n.left = child
	}
	// n.Update()
}

func isRed[T constraints.Ordered](root *rbTreeNode[T]) bool {
	return root != nil && root.red()
}
