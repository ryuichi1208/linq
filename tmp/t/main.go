package main

import "fmt"

type P1 struct {
	p []*P2
}

type P2 struct {
	age int
}

func ch(pp *P1) {
	for i, _ := range pp.p {
		pp.p[i].age = 2
	}
}

func ch2(pp P1) {
	for i, _ := range pp.p {
		pp.p[i].age = 3
	}
}

func main() {
	p2 := &P2{1}
	p3 := &P2{2}

	var P2S []*P2
	P2S = append(P2S, p2, p3)
	p1 := P1{P2S}

	ch(&p1)
	for i, v := range p1.p {
		fmt.Println(i, v.age)
	}

	ch2(p1)
	for i, v := range p1.p {
		fmt.Println(i, v.age)
	}

	n := &P2{10}
	fmt.Println((*n).age)

	a := []int{1, 3, 5, 6}
	fmt.Println(a)
}
