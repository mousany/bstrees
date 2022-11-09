package node

import (
	"bstrees/pkg/type/number"
	"math/rand"
)

type TreapNode[T number.Number] struct {
	Value  T
	Left   *TreapNode[T]
	Right  *TreapNode[T]
	Weight uint32
	Size   uint32
}

func New[T number.Number](value T) *TreapNode[T] {
	return &TreapNode[T]{Value: value, Left: nil, Right: nil, Weight: rand.Uint32(), Size: 1}
}

func (this *TreapNode[T]) Update() {
	this.Size = 1
	if this.Left != nil {
		this.Size += this.Left.Size
	}
	if this.Right != nil {
		this.Size += this.Right.Size
	}
}
