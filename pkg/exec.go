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

// Map executes `taskFunc` concurrently for each item in `list`
func Map(list interface{}, taskFunc interface{}, options *Options) chan *ReturnValue {
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

	// execute `taskFunc` for each item in the list
	for i := 0; i < nTasks; i++ {
		elem := actualList.Index(i)
		go func(elem reflect.Value) {
			defer wg.Done()
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
		}(elem)
	}
	return out
}
