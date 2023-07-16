package splay

import (
	"github.com/yanglinshu/bstrees/v2"
	"golang.org/x/exp/constraints"
)

type Splay[T constraints.Ordered] struct {
	superRoot *splayNode[T]
}

func (t *Splay[T]) root() *splayNode[T] {
	return t.superRoot.right
}

func (t *Splay[T]) setRoot(root *splayNode[T]) {
	t.superRoot.setChild(root, true)
}

func New[T constraints.Ordered]() *Splay[T] {
	return &Splay[T]{
		superRoot: newSplayNode(T(rune(0))),
	}
}

func search[T constraints.Ordered](root *splayNode[T], value T) *splayNode[T] {
	for p := root; p != nil; {
		if p.value == value {
			return p
		} else if value < p.value {
			p = p.left
		} else {
			p = p.right
		}
	}
	return nil
}

func at[T constraints.Ordered](root *splayNode[T], k uint) *splayNode[T] {
	for p := root; p != nil; {
		leftSize := uint(0)
		if p.left != nil {
			leftSize = p.left.size
		}
		if leftSize < k && leftSize+p.rec >= k {
			// SplayRotate(p, root)
			return p
		} else if leftSize+p.rec < k {
			k -= leftSize + p.rec
			p = p.right
		} else {
			p = p.left
		}
	}
	return nil
}

func insert[T constraints.Ordered](root *splayNode[T], value T) *splayNode[T] {
	if root == nil {
		return newSplayNode(value)
	} else {
		superRoot := root.parent

		for p := root; p != nil; {
			p.size += 1
			if value == p.value {
				p.rec += 1
				splayRotate(p, root)
				break
			} else if value < p.value {
				if p.left == nil {
					p.setChild(newSplayNode(value), false)
					splayRotate(p.left, root)
					break
				} else {
					p = p.left
				}
			} else {
				if p.right == nil {
					p.setChild(newSplayNode(value), true)
					splayRotate(p.right, root)
					break
				} else {
					p = p.right
				}
			}
		}

		return superRoot.right
	}
}

func delete[T constraints.Ordered](root *splayNode[T], value T) *splayNode[T] {
	if root == nil {
		return nil
	}
	superRoot := root.parent
	p := search(root, value)
	if p == nil {
		return root
	}
	splayRotate(p, root)
	if p.rec > 1 {
		p.rec -= 1
		p.size -= 1
	} else {
		if p.left == nil && p.right == nil {
			superRoot.setChild(nil, true)
		} else if p.left == nil {
			superRoot.setChild(p.right, true)
		} else if p.right == nil {
			superRoot.setChild(p.left, true)
		} else {
			maxLeft := p.left
			for maxLeft.right != nil {
				maxLeft.size -= 1
				maxLeft = maxLeft.right
			}
			splayRotate(maxLeft, superRoot.right)
			maxLeft.setChild(p.right, true)
			superRoot.setChild(maxLeft, true)
			superRoot.right.update()
		}
	}

	return superRoot.right
}

func (t *Splay[T]) Insert(value T) {
	t.setRoot(insert(t.root(), value))
}

func (t *Splay[T]) Delete(value T) {
	t.setRoot(delete(t.root(), value))
}

func (t *Splay[T]) Contains(value T) bool {
	return search(t.root(), value) != nil
}

func (t *Splay[T]) At(k uint) (T, error) {
	result := at(t.root(), k)
	if result == nil {
		return T(rune(0)), bstrees.ErrIndexIsOutOfRange
	}
	return result.value, nil
}

func (t *Splay[T]) Size() uint {
	if t.root() == nil {
		return 0
	}
	return t.root().size
}

func (t *Splay[T]) Empty() bool {
	return t.root() == nil
}

func (t *Splay[T]) Clear() {
	t.setRoot(nil)
}

func (t *Splay[T]) Index(value T) uint {
	// return Rank(t.Root, value)
	p := search(t.root(), value)
	if p == nil {
		prev := predecessor(t.root(), value)
		if prev != nil {
			splayRotate(prev, t.root())
			if prev.left != nil {
				return prev.left.size + prev.rec + 1
			}
			return prev.rec + 1
		}
		return 1
	}
	splayRotate(p, t.root())
	if p.left != nil {
		return p.left.size + 1
	}
	return 1
}

func predecessor[T constraints.Ordered](root *splayNode[T], value T) *splayNode[T] {
	var result *splayNode[T]
	for p := root; p != nil; {
		if value > p.value {
			result = p
			p = p.right
		} else {
			p = p.left
		}
	}
	return result
}

func (t *Splay[T]) Predecessor(value T) (T, error) {
	prev := predecessor(t.root(), value)
	if prev == nil {
		return T(rune(0)), bstrees.ErrPredecessorDoesNotExist
	}
	return prev.value, nil
}

func successor[T constraints.Ordered](root *splayNode[T], value T) *splayNode[T] {
	var result *splayNode[T]
	for p := root; p != nil; {
		if value < p.value {
			result = p
			p = p.left
		} else {
			p = p.right
		}
	}
	return result
}

func (t *Splay[T]) Successor(value T) (T, error) {
	next := successor(t.root(), value)
	if next == nil {
		return T(rune(0)), bstrees.ErrSuccessorDoesNotExist
	}
	return next.value, nil
}
