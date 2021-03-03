package main

import (
	"fmt"
	"strconv"
)

func f(l []int) []string {
	var s []string
	fmt.Println(&l[0])
	for _, v := range l {
		if m := v * v; m%2 == 0 {
			fmt.Println(v)
			s = append(s, strconv.Itoa(v))
		}
	}
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	return s
}

func main() {
	var l []int
	l = append(l, 1, 2, 3, 4)
	fmt.Println(&l[0])
	fmt.Println(f(l))
}
