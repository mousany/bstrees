package rbtree

import "bstrees/internal/node"

type RBColor bool

const (
	Red   RBColor = true
	Black RBColor = false
)

type rbTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]
	color RBColor
}

func newRBTreeNode[T node.Ordered](value T) *rbTreeNode[T] {
	n := &rbTreeNode[T]{node.New(value), Red}
	n.SetLeft((*rbTreeNode[T])(nil))
	n.SetRight((*rbTreeNode[T])(nil))
	return n
}

func (n *rbTreeNode[T]) IsNil() bool {
	return n == nil
}

func (n *rbTreeNode[T]) Color() RBColor {
	return n.color
}

func (n *rbTreeNode[T]) SetColor(color RBColor) {
	n.color = color
}

func (n *rbTreeNode[T]) IsRed() bool {
	return n != nil && n.color == Red
}

func (n *rbTreeNode[T]) IsBlack() bool {
	return n == nil || n.color == Black
}
