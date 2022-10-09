package main

import (
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := Api{"192.168.1.92", "8080"}
	go api.start(&wg)
	wg.Wait()
}
