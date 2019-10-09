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

func Cube(i int) int {
	s := i * i * i
	return s
}

func main() {

	list := []int{1, 2, 3, 4, 5, 6}
	res := pkg.Map(list, Square, &pkg.Options{Concurrency: 1})

	for i := range res {
		fmt.Println(i.Result.(int))
	}

	in := make(chan interface{}, 10)
	for _, i := range list {
		in <- i
	}
	close(in)
	res2 := pkg.Pipe(
		pkg.Pipe(in, Cube, &pkg.Options{Concurrency: 1}),
		Square,
		&pkg.Options{Concurrency: 1},
	)
	for i := range res2 {
		fmt.Println(i.(int))
	}

	pkg.Concurrent(pkg.Options{}, func(done pkg.ResultChan) {
		fmt.Println("Task 1")
		done <- 1
	}, func(done pkg.ResultChan) {
		fmt.Println("Task 2")
		done <- 1
	}, func(done pkg.ResultChan) {
		fmt.Println("Task 3")
		done <- 1
	})

	fmt.Println("THE END")
}
