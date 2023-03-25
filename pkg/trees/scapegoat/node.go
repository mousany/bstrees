package scapegoat

import "bstrees/internal/order"

type nodeState bool

const (
	inactive nodeState = false
	active   nodeState = true
)

type scapeGoatTreeNode[T order.Number] struct {
	value  T
	left   *scapeGoatTreeNode[T]
	right  *scapeGoatTreeNode[T]
	state  nodeState
	size   uint // Number of active nodes in the subtree
	weight uint // Number of nodes in the subtree
}

func newScapeGoatTreeNode[T order.Number](value T) *scapeGoatTreeNode[T] {
	return &scapeGoatTreeNode[T]{
		value:  value,
		left:   nil,
		right:  nil,
		state:  active,
		size:   1,
		weight: 1,
	}
}

func (n *scapeGoatTreeNode[T]) inactive() bool {
	return n.state == inactive
}

func (n *scapeGoatTreeNode[T]) active() bool {
	return n.state == active
}

func (n *scapeGoatTreeNode[T]) update() {
	if n.active() {
		n.size = 1
	} else {
		n.size = 0
	}
	n.weight = 1
	if n.left != nil {
		n.size += n.left.size
		n.weight += n.left.weight
	}
	if n.right != nil {
		n.size += n.right.size
		n.weight += n.right.weight
	}
}

func (n *scapeGoatTreeNode[T]) deactivate() {
	n.state = inactive
	n.update()
}
