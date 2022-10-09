package main

import (
	"sync"
)

//Resp - Generic response
type Resp struct {
	Result string `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := Api{"192.168.1.92", "8080"}
	go api.start(&wg)
	wg.Wait()
}
