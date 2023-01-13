package splay

import "bstrees/internal/node"

type splayTreeNode[T node.Ordered] struct {
	*node.BaseTreeNode[T]

	parent *splayTreeNode[T]
	rec    uint // This field is Splay only
	// Because Splay operation will scatter nodes with the same value
	// While traditional BST search mechanics is too slow on Splay
}

func newSplayTreeNode[T node.Ordered](value T) *splayTreeNode[T] {
	n := &splayTreeNode[T]{
		BaseTreeNode: node.New(value),
		parent:       (*splayTreeNode[T])(nil),
		rec:          1,
	}
	n.SetLeft((*splayTreeNode[T])(nil))
	n.SetRight((*splayTreeNode[T])(nil))
	return n
}

func (n *splayTreeNode[T]) IsNil() bool {
	return n == nil
}

func (n *splayTreeNode[T]) Parent() *splayTreeNode[T] {
	return n.parent
}

func (n *splayTreeNode[T]) SetParent(parent *splayTreeNode[T]) {
	n.parent = parent
}

func (n *splayTreeNode[T]) Rec() uint {
	return n.rec
}

func (n *splayTreeNode[T]) SetRec(rec uint) {
	n.rec = rec
}

func (n *splayTreeNode[T]) SetLeft(left node.Noded[T]) {
	n.BaseTreeNode.SetLeft(left)
	if !left.IsNil() {
		left.(*splayTreeNode[T]).SetParent(n)
	}
}

func (n *splayTreeNode[T]) SetRight(right node.Noded[T]) {
	n.BaseTreeNode.SetRight(right)
	if !right.IsNil() {
		right.(*splayTreeNode[T]).SetParent(n)
	}
}

func (n *splayTreeNode[T]) SetChild(right bool, child node.Noded[T]) {
	if right {
		n.SetRight(child)
	} else {
		n.SetLeft(child)
	}
}

func (n *splayTreeNode[T]) Update() {
	n.BaseTreeNode.Update()
	n.SetSize(n.Size() + n.Rec() - 1)
}
