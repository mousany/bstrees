package treap

import (
	"bstrees/internal/node"
	"math/rand"
)

type treapTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]
	weight uint // Random weight
}

func newTreapTreeNode[T node.Ordered](value T) *treapTreeNode[T] {
	n := &treapTreeNode[T]{BaseTreeNode: node.New(value), weight: uint(rand.Uint32())}
	n.SetLeft((*treapTreeNode[T])(nil))
	n.SetRight((*treapTreeNode[T])(nil))
	return n
}

func (n *treapTreeNode[T]) Weight() uint {
	return n.weight
}

func (n *treapTreeNode[T]) IsNil() bool {
	return n == nil
}
