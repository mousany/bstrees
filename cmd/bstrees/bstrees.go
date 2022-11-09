package main

import (
	"bstrees/pkg/treap/vanilla"
	"bstrees/pkg/util"
	"bufio"
	"fmt"
	"os"
)

func main() {
	tree := vanilla.New[int]()
	gin := bufio.NewReader(os.Stdin)
	n, err := util.Read[int](gin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		opt, err := util.Read[int](gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		value, err := util.Read[int](gin)
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
