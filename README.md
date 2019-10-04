# chango
A go package for running tasks concurrently and control concurreny, task completion and errors.

# Examples
```go
func main() {

	list := []int{1, 2, 3, 4, 5, 6}
	res := chango.Map(list, Square, &pkg.Options{Concurrency: 10})

	for i := range res {
		fmt.Println(i.Result.(int))
	}

}
```