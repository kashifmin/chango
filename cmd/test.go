package main

import (
	"fmt"

	"github.com/kashifmin/chango/pkg"
)

type ArgsList interface {
	Len() int
}

func Square(i int) int {
	s := i * i
	return s
}

func main() {

	list := []int{1, 2, 3, 4, 5, 6}
	res := pkg.Map(list, Square, &pkg.Options{Concurrency: 10})

	for i := range res {
		fmt.Println(i.Result.(int))
	}

}
