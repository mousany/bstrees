package rb

import (
	"github.com/yanglinshu/bstrees/v2/internal/order"
)

func singleRotate[T order.Ordered](root *rbTreeNode[T], direction bool) *rbTreeNode[T] {
	save := root.child(!direction)
	root.setChild(!direction, save.child(direction))
	save.setChild(direction, root)
	root.Update()
	save.Update()
	root.color = red
	save.color = black
	return save
}

func doubleRotate[T order.Ordered](root *rbTreeNode[T], direction bool) *rbTreeNode[T] {
	root.setChild(!direction, singleRotate(root.child(!direction), !direction))
	return singleRotate(root, direction)
}
