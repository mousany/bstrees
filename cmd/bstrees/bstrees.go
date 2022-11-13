package main

import (
	"bstrees/pkg/avltree"
	"bstrees/pkg/util/console"
	"bufio"
	"fmt"
	"os"
)

func main() {
	tree := avltree.New[int]()
	gin := bufio.NewReader(os.Stdin)
	n, err := console.Read[int](gin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		opt, err := console.Read[int](gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		value, err := console.Read[int](gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		switch opt {
		case 1:
			tree.Insert(value)
		case 2:
			tree.Delete(value)
		case 3:
			fmt.Println(tree.Rank(value))
		case 4:
			{
				kth, err := tree.Kth(uint32(value))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		case 5:
			{
				kth, err := tree.Prev(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		case 6:
			{
				kth, err := tree.Next(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		}
	}
}
