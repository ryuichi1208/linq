package main

import (
	"fmt"
	"time"
)

type List struct {
	node int
}

type Human struct {
	name string
	age  int
	list *[]List
}

func ListChanged(l *[]List) {
	for k, v := range *l {
		(*l)[k].node = 20
		fmt.Println(k, v, (*l)[k])
	}
}

func (h *Human) New() {
	// ListChanged(h.list)

	for i, n := range *h.list {
		fmt.Println("aaa", i, n)
		(*h.list)[i].node = 20
		fmt.Println((*h.list)[i].node)
	}
}

func p(h map[string]Human) {
	for k, v := range h {
		v.name = "yasushi"
		fmt.Println(k, v.name)
		for _, _ = range *v.list {
			// i.node = 20
			ll := h[k].list
			// j.node = 20
			// fmt.Println(*ll, i, j.node)
			for m, n := range *ll {
				n.node = 20
				fmt.Println(m, n.node)
			}
		}
	}
}

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	// time.Sleep(1600 * time.Millisecond)
	// ticker.Stop()
	// fmt.Println("Ticker stopped")

	h := make(map[string]Human)
	l := List{10}
	ls := &[]List{l, l}

	h["test"] = Human{"takashi", 14, ls}

	// p(h)
	h2 := Human{"yamashi", 15, ls}
	fmt.Println(h2)
	h2.New()
	for k, v := range *h2.list {
		fmt.Println(k, v.node)
	}
}
