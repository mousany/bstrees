package splay

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
	"bstrees/pkg/errors"
)

type SplayTree[T node.Ordered] struct {
	*tree.BaseTree[T] // In splay tree, root is not the root of the tree, but the right child of superRoot
}

func New[T node.Ordered]() *SplayTree[T] {
	tr := &SplayTree[T]{tree.New[T]()}
	tr.BaseTree.SetRoot(newSplayTreeNode(T(rune(0))))
	return tr
}

func (tr *SplayTree[T]) Root() node.Noded[T] {
	return tr.BaseTree.Root().Right()
}

func (tr *SplayTree[T]) SetRoot(root node.Noded[T]) {
	tr.BaseTree.Root().SetRight(root)
}

func insert[T node.Ordered](root *splayTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return newSplayTreeNode(value)
	} else {
		superRoot := root.Parent()
		for p := root; p != nil; {
			p.SetSize(p.Size() + 1)
			if value == p.Value() {
				p.SetRec(p.Rec() + 1)
				splayRotate(p, root)
				break
			} else if value < p.Value() {
				if p.Left().IsNil() {
					p.SetLeft(newSplayTreeNode(value))
					splayRotate(p.Left().(*splayTreeNode[T]), root)
					break
				} else {
					p = p.Left().(*splayTreeNode[T])
				}
			} else {
				if p.Right().IsNil() {
					p.SetRight(newSplayTreeNode(value))
					splayRotate(p.Right().(*splayTreeNode[T]), root)
					break
				} else {
					p = p.Right().(*splayTreeNode[T])
				}
			}
		}
		return superRoot.Right()
	}
}

func (tr *SplayTree[T]) Insert(value T) {
	tr.SetRoot(insert(tr.Root().(*splayTreeNode[T]), value))
}

func Delete[T node.Ordered](root *splayTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return (*splayTreeNode[T])(nil)
	}
	superRoot := root.Parent()
	p := tree.Find(node.Noded[T](root), value).(*splayTreeNode[T])
	if p == nil {
		return root
	}
	splayRotate(p, root)
	if p.Rec() > 1 {
		p.SetRec(p.Rec() - 1)
		p.SetSize(p.Size() - 1)
	} else {
		if p.Left().IsNil() && p.Right().IsNil() {
			superRoot.SetRight((*splayTreeNode[T])(nil))
		} else if p.Left().IsNil() {
			superRoot.SetRight(p.Right())
		} else if p.Right().IsNil() {
			superRoot.SetRight(p.Left())
		} else {
			maxLeft := p.Left()
			for !maxLeft.Right().IsNil() {
				maxLeft.SetSize(maxLeft.Size() - 1)
				maxLeft = maxLeft.Right()
			}
			splayRotate(maxLeft.(*splayTreeNode[T]), superRoot.Right().(*splayTreeNode[T]))
			maxLeft.SetRight(p.Right())
			superRoot.SetRight(maxLeft)
			superRoot.Right().Update()
		}
	}

	return superRoot.Right()
}

func (tr *SplayTree[T]) Delete(value T) {
	tr.SetRoot(Delete(tr.Root().(*splayTreeNode[T]), value))
}

func kth[T node.Ordered](root *splayTreeNode[T], k uint) node.Noded[T] {
	for p := root; p != nil; {
		leftSize := uint(0)
		if !p.Left().IsNil() {
			leftSize = p.Left().Size()
		}
		if leftSize < k && leftSize+p.Rec() >= k {
			splayRotate(p, root)
			return p
		} else if leftSize+p.Rec() < k {
			k -= leftSize + p.Rec()
			p = p.Right().(*splayTreeNode[T])
		} else {
			p = p.Left().(*splayTreeNode[T])
		}
	}
	return (*splayTreeNode[T])(nil)
}

func (tr *SplayTree[T]) Kth(k uint) (T, error) {
	result := kth(tr.Root().(*splayTreeNode[T]), k)
	if result.IsNil() {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.Value(), nil
}

func (tr *SplayTree[T]) Rank(value T) uint {
	p := tree.Find(tr.Root(), value).(*splayTreeNode[T])
	if p == nil {
		prev := tree.Prev(tr.Root(), value).(*splayTreeNode[T])
		if prev != nil {
			splayRotate(prev, tr.Root().(*splayTreeNode[T]))
			if !prev.Left().IsNil() {
				return prev.Left().Size() + prev.Rec() + 1
			}
			return prev.Rec() + 1
		}
		return 1
	}
	splayRotate(p, tr.Root().(*splayTreeNode[T]))
	if !p.Left().IsNil() {
		return p.Left().Size() + 1
	}
	return 1
}

func (tr *SplayTree[T]) Size() uint {
	if tr.Root().IsNil() {
		return 0
	}
	return tr.Root().Size()
}

func (tr *SplayTree[T]) Empty() bool {
	return tr.Root().IsNil()
}

func (tr *SplayTree[T]) Clear() {
	tr.SetRoot((*splayTreeNode[T])(nil))
}

func (tr *SplayTree[T]) Contains(value T) bool {
	return tree.Find(tr.Root(), value) != nil
}

func (tr *SplayTree[T]) Prev(value T) (T, error) {
	result := tree.Prev(tr.Root(), value)
	if tree.IsNil(result) {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return result.Value(), nil
}

func (tr *SplayTree[T]) Next(value T) (T, error) {
	result := tree.Next(tr.Root(), value)
	if tree.IsNil(result) {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return result.Value(), nil
}

func (tr *SplayTree[T]) String() string {
	return tree.String(tr.Root())
}
