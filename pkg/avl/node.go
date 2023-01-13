package avltree

import (
	"bstrees/internal/math"
	"bstrees/internal/node"
)

type avlTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]     // Embedded node
	height                int // Height of the node
}

func newAvlTreeNode[T node.Ordered](value T) *avlTreeNode[T] {
	n := &avlTreeNode[T]{node.New(value), 0}
	n.SetLeft((*avlTreeNode[T])(nil))
	n.SetRight((*avlTreeNode[T])(nil))
	return n
}

func (n *avlTreeNode[T]) Height() int {
	if n == nil {
		return -1
	}
	return n.height
}

func (n *avlTreeNode[T]) SetHeight(height int) {
	n.height = height
}

func (n *avlTreeNode[T]) IsNil() bool {
	return n == nil
}

func (n *avlTreeNode[T]) Update() {
	n.BaseTreeNode.Update()
	n.height = 0
	if !n.Left().IsNil() {
		n.height = math.Max(n.height, n.Left().(*avlTreeNode[T]).Height()+1)
	}
	if !n.Right().IsNil() {
		n.height = math.Max(n.height, n.Right().(*avlTreeNode[T]).Height()+1)
	}
}
