package anderson

import "bstrees/internal/node"

type andersonTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]
	level uint
}

func newAndersonTreeNode[T node.Ordered](value T, level uint) *andersonTreeNode[T] {
	n := &andersonTreeNode[T]{
		BaseTreeNode: node.New(value),
		level:        level,
	}
	n.SetLeft((*andersonTreeNode[T])(nil))
	n.SetRight((*andersonTreeNode[T])(nil))
	return n
}

func (n *andersonTreeNode[T]) IsNil() bool {
	return n == nil
}

func (n *andersonTreeNode[T]) Level() uint {
	return n.level
}

func (n *andersonTreeNode[T]) SetLevel(level uint) {
	n.level = level
}
