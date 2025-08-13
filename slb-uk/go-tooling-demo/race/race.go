package main

import (
	"fmt"
	"sync"
)

func main() {
	var x int
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			x++ // intentional race
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("x:", x)
}
