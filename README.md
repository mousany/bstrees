# BSTrees: Implementation of Binary Search Tree algorithms in Go

This repository contains implementation of Binary Search Trees algorithms in Go. Those trees are tested on the [template problem](https://www.luogu.com.cn/problem/P3369) and the [enhanced template problem](https://www.luogu.com.cn/problem/P6136) of [Luogu](https://www.luogu.com.cn/).

## Usage
All the trees are implemented in the `bstree` package. The `bstree` package contains the following trees:
- `bstree.avl.AVLTree`: AVL Tree
- `bstree.rb.RBTree`: Red-Black Tree
- `bstree.anderson.AndersonTree`: Anderson Tree
- `bstree.treap.Treap`: Treap
- `bstree.fhq.FHQTreap`: FHQ Rotateless Treap
- `bstree.splay.Splay`: Splay Tree
- `bstree.scapegoat.ScapegoatTree`: Scapegoat Tree

These trees are implemented in the same way and share an uniform interface. Here is an example of using the `bstree.avl.AVLTree`:
```go
package main

import (
    "fmt"
    "github.com/yanglinshu/bstrees/avl"
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
    tree.Index(6) // Output: 5
    tree.At(5) // Output: 6
    tree.Predecessor(6) // Output: 4
    tree.Successor(6) // Output: 7
}
```

## Production
It might be better to try bstrees out on a hobby project first. Bstrees does not aim to be a production-ready library. It is migrated from some ACM contest code and is still having performance issues. And there is not guaranteed to be bug-free and the API might change in the future. However, it will be a good choice for you to learn about binary search trees.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.