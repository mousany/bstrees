package main

import (
	"bstrees/pkg/splay"
	"bufio"
	"fmt"
	"os"
)

func Read(istream *bufio.Reader) (int, error) {
	res, sign := int(0), 1
	readed := false
	c, err := istream.ReadByte()
	for ; err == nil && (c < '0' || c > '9'); c, err = istream.ReadByte() {
		if c == '-' {
			sign = -1
		}
	}
	for ; err == nil && c >= '0' && c <= '9'; c, err = istream.ReadByte() {
		readed = true
		res = res<<3 + res<<1 + int(c-'0')
	}
	if !readed {
		return 0, err
	}
	return res * int(sign), err
}

func ReadWithPanic(gin *bufio.Reader) int {
	value, err := Read(gin)
	if err != nil {
		panic(err)
	}
	return value
}

const Debug = false

func main() {

	lastSize := uint(0)
	tr := splay.New[int]()
	gin := bufio.NewReader(os.Stdin)
	n := ReadWithPanic(gin)
	for i := 0; i < n; i++ {
		opt := ReadWithPanic(gin)
		value := ReadWithPanic(gin)
		if Debug {
			fmt.Println("----------------")
			fmt.Println(opt, value)
		}
		switch opt {
		case 1:
			tr.Insert(value)
		case 2:
			tr.Delete(value)
		case 3:
			fmt.Println(tr.Rank(value))
		case 4:
			{
				kth, err := tr.Kth(uint(value))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		case 5:
			{
				kth, err := tr.Prev(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		case 6:
			{
				kth, err := tr.Next(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		}
		if Debug && (opt == 1 || opt == 2) {
			err := splay.CheckBSTProperty(tr.Root())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if opt == 1 {
				if tr.Size() != lastSize+1 {
					fmt.Println("Insert error")
					os.Exit(1)
				}
			} else {
				if tr.Size() != lastSize-1 {
					fmt.Println("Delete error")
					os.Exit(1)
				}
			}
			lastSize = tr.Size()
			fmt.Println(tr.String())
		}
	}
}
