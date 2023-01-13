package node

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type Number interface {
	Integer | Float
}

type Ordered interface {
	Number | ~string
}

type Valued[T Ordered] interface {
	Value() T
	SetValue(T)
}

type Sized[T Ordered] interface {
	Size() uint
	SetSize(uint)
}

type Binary[T Ordered] interface {
	Left() Noded[T]
	SetLeft(Noded[T])
	Right() Noded[T]
	SetRight(Noded[T])
	Child(bool) Noded[T]
	SetChild(bool, Noded[T])
}

type Noded[T Ordered] interface {
	Valued[T]
	Sized[T]
	Binary[T]
	Update()
	IsNil() bool
}

type BaseTreeNode[T Ordered] struct {
	value T
	size  uint
	left  Noded[T]
	right Noded[T]
}

func New[T Ordered](value T) *BaseTreeNode[T] {
	return &BaseTreeNode[T]{value, 1, nil, nil}
}

func (n *BaseTreeNode[T]) Value() T {
	return n.value
}

func (n *BaseTreeNode[T]) SetValue(value T) {
	n.value = value
}

func (n *BaseTreeNode[T]) Size() uint {
	return n.size
}

func (n *BaseTreeNode[T]) SetSize(size uint) {
	n.size = size
}

func (n *BaseTreeNode[T]) Left() Noded[T] {
	return n.left
}

func (n *BaseTreeNode[T]) SetLeft(left Noded[T]) {
	n.left = left
}

func (n *BaseTreeNode[T]) Right() Noded[T] {
	return n.right
}

func (n *BaseTreeNode[T]) SetRight(right Noded[T]) {
	n.right = right
}

func (n *BaseTreeNode[T]) Child(right bool) Noded[T] {
	if right {
		return n.right
	}
	return n.left
}

func (n *BaseTreeNode[T]) SetChild(right bool, child Noded[T]) {
	if right {
		n.right = child
	} else {
		n.left = child
	}
}

func (n *BaseTreeNode[T]) IsNil() bool {
	return n == nil
}

func (n *BaseTreeNode[T]) Update() {
	n.size = 1
	if !n.left.IsNil() {
		n.size += n.left.Size()
	}
	if !n.right.IsNil() {
		n.size += n.right.Size()
	}
}
