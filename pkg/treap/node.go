package treap

import (
	"bstrees/internal/node"
	"math/rand"
)

type TreapTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]
	weight uint // Random weight
}

func NewTreapTreeNode[T node.Ordered](value T) *TreapTreeNode[T] {
	n := &TreapTreeNode[T]{BaseTreeNode: node.New(value), weight: uint(rand.Uint32())}
	n.SetLeft((*TreapTreeNode[T])(nil))
	n.SetRight((*TreapTreeNode[T])(nil))
	return n
}

func (n *TreapTreeNode[T]) Weight() uint {
	return n.weight
}

func (n *TreapTreeNode[T]) IsNil() bool {
	return n == nil
}
