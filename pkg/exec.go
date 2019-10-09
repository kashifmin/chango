package pkg

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Options struct {
	Concurrency int
}

type ReturnValue struct {
	Error  error
	Result interface{}
}

type ResultChan chan interface{}

// Map executes `taskFunc` concurrently for each item in `list`
func Map(list interface{}, taskFunc interface{}, options *Options) <-chan *ReturnValue {
	// check for correct types by reflection
	actualFunc := reflect.Indirect(reflect.ValueOf(taskFunc))
	actualList := reflect.ValueOf(list)
	if actualFunc.Kind() != reflect.Func {
		panic(errors.New("taskFunc must be a function"))
	}
	if actualList.Kind() != reflect.Slice && actualList.Kind() != reflect.Array {
		fmt.Println(actualList.Kind())
		panic(errors.New("List Must be an array or slice"))
	}
	if options.Concurrency < 1 {
		panic(errors.New("Concurrency is less than 1! You do not want to be deadlocked."))
	}

	// channel for sending the results of completed functions
	out := make(chan *ReturnValue)

	nTasks := actualList.Len()

	// Setup wait group to control execution of all tasks
	var wg sync.WaitGroup
	wg.Add(nTasks)
	go func() {
		wg.Wait()
		close(out)
	}()

	sem := make(chan int, options.Concurrency)

	// execute `taskFunc` for each item in the list
	for i := 0; i < nTasks; i++ {
		elem := actualList.Index(i)
		go func(elem reflect.Value) {
			defer wg.Done()
			sem <- 1
			result := actualFunc.Call([]reflect.Value{elem})

			// wrap function results as `ReturnValue` type
			retVal := &ReturnValue{}
			n := len(result)
			if n == 1 {
				retVal.Result = result[0].Interface()
			} else if n == 2 {
				retVal.Result = result[0].Interface()
				retVal.Error = result[1].Interface().(error)
			}
			out <- retVal
			<-sem
		}(elem)
	}
	return out
}

// Pipe reads from channel `in` and applies `taskFunc` on it
// results are sent out using the channel returned
func Pipe(in <-chan interface{}, taskFunc interface{}, options *Options) <-chan interface{} {
	actualFunc := reflect.Indirect(reflect.ValueOf(taskFunc))
	if actualFunc.Kind() != reflect.Func {
		panic(errors.New("taskFunc must be a function"))
	}

	// channel for sending the results of completed functions
	out := make(chan interface{})

	var wg sync.WaitGroup

	for elem := range in {
		wg.Add(1)
		go func(e interface{}) {
			defer wg.Done()

			result := actualFunc.Call([]reflect.Value{reflect.ValueOf(e)})

			// wrap function results as `ReturnValue` type
			var retVal interface{}
			n := len(result)
			if n == 1 {
				retVal = result[0].Interface()
			} else {
				panic(errors.New("taskFunc returned more than one value!"))
			}

			out <- retVal

		}(elem)

	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Concurrent(options Options, taskFunc ...func(ResultChan)) {
	n := len(taskFunc)

	doneChan := make(ResultChan, n)

	for _, task := range taskFunc {
		go task(doneChan)
	}

	for i := 0; i < n; i++ {
		<-doneChan
	}
	return
}
