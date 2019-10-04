# chango
A go package for running tasks concurrently and control concurreny, task completion and errors.

# Examples

- Apply a function to items of a list concurrently
```go
func main() {

	list := []int{1, 2, 3, 4, 5, 6}
	res := chango.Map(list, Square, &pkg.Options{Concurrency: 10})

	for i := range res {
		fmt.Println(i.Result.(int))
	}

}
```

- Apply a function to values recived from a channel
```go
func main() {

	list := []int{1, 2, 3, 4, 5, 6}

	in := make(chan interface{}, 10)
	for _, i := range list {
		in <- i
	}
	close(in)
	res2 := pkg.Pipe(
		pkg.Pipe(in, Cube, &pkg.Options{Concurrency: 1}), 
		Square, 
		&pkg.Options{Concurrency: 1}
	)
	for i := range res2 {
		fmt.Println(i.(int))
	}
}
```
You can combine `Pipe` do some powerful stuff. 
For example, you can do Map, Filter on the items in the channel concurrently!
NOTE: Remember to close the input channel once all items are sent, or else it will result in a deadlock.