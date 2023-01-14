# BSTrees: Go implementation of some Binary Search Trees

This repository contains the implementation of some Binary Search Trees in Go. Those trees are tested on the [template problem](https://www.luogu.com.cn/problem/P3369) and the [enhanced template problem](https://www.luogu.com.cn/problem/P6136) of [Luogu](https://www.luogu.com.cn/).

## Usage
All the trees are implemented in the `bstree` package. The `bstree` package contains the following trees:
- `bstree.avl.AVLTree`: AVL Tree
- `bstree.rb.RBTree`: Red-Black Tree
- `bstree.anderson.AndersonTree`: Anderson Tree
- `bstree.treap.TreapTree`: Treap
- `bstree.fhq.FHQTree`: FHQ Rotateless Treap
- `bstree.splay.SplayTree`: Splay Tree (Removed due to the lack of a good implementation)
- `bstree.scapegoat.ScapegoatTree`: Scapegoat Tree (Removed due to the lack of a good implementation)

These trees are implemented in the same way and share an uniform interface. Here is an example of using the `bstree.avl.AVLTree`:
```go
package main

import (
    "fmt"
    "github.com/yanglinshu/bstrees/pkg/avl"
)

func main() {
    tree := avl.New()
    tree.Insert(1)
    tree.Insert(2)
    tree.Insert(3)
    tree.Insert(4)
    tree.Insert(5)
    tree.Insert(6)
    tree.Insert(7)
    tree.Insert(8)
    tree.Insert(9)
    tree.Insert(10)
    tree.Delete(5)
    tree.Rank(6) // Output: 5
    tree.Kth(5) // Output: 6
    tree.Prev(6) // Output: 4
    tree.Next(6) // Output: 7
}
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.