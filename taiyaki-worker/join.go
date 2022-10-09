package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func NodeJoin(workerIP string, workerPort string, serverIP string, serverPort string, joinKey string) string {
	fmt.Println("Joining the cluster...")

	reqJoin :=  NodeJoinReq{NodeIP: workerIP, NodePort: workerPort, JoinKey: joinKey}
	reqBuffer := new(bytes.Buffer)
	json.NewEncoder(reqBuffer).Encode(reqJoin)

	resp, err := reqServer("http://"+serverIP+":8080/node/join", reqBuffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp))

	// Complete the error handling part here
	fmt.Println("Joined the cluster")
	return ""
}
