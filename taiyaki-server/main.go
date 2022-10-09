package main

import (
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := Api{"127.0.0.1", "8080"}
	go api.start(&wg)
	wg.Wait()
}