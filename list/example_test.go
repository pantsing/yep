package list_test

import (
	"fmt"

	"tsing.cn/yep/list"
)

func ExampleFIFO() {
	// Create a new list and put some numbers in it.
	l, err := list.New(list.FIFO, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}

	l.Push(list.StringElem("1"))
	l.Push(list.StringElem("2"))
	l.Push(list.StringElem("3"))
	l.Push(list.StringElem("4"))
	for i := 0; i < 4; i++ {
		e, err := l.Pop()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(e)
		}
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleLIFO() {
	// Create a new list and put some numbers in it.
	l, err := list.New(list.LIFO, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}

	l.Push(list.StringElem("1"))
	l.Push(list.StringElem("2"))
	l.Push(list.StringElem("3"))
	l.Push(list.StringElem("4"))
	for i := 0; i < 4; i++ {
		e, err := l.Pop()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(e)
		}
	}
	// Output:
	// 4
	// 3
	// 2
	// 1
}
