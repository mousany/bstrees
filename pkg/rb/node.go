package rb

import "bstrees/internal/node"

type rbColor bool

const (
	red   rbColor = true
	black rbColor = false
)

type rbTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]
	color rbColor
}

func newRBTreeNode[T node.Ordered](value T) *rbTreeNode[T] {
	n := &rbTreeNode[T]{node.New(value), red}
	n.SetLeft((*rbTreeNode[T])(nil))
	n.SetRight((*rbTreeNode[T])(nil))
	return n
}

func (n *rbTreeNode[T]) IsNil() bool {
	return n == nil
}

func (n *rbTreeNode[T]) Color() rbColor {
	return n.color
}

func (n *rbTreeNode[T]) SetColor(color rbColor) {
	n.color = color
}

func (n *rbTreeNode[T]) IsRed() bool {
	return n != nil && n.color == red
}

func (n *rbTreeNode[T]) IsBlack() bool {
	return n == nil || n.color == black
}
