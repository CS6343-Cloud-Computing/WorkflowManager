package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func NodeJoin(workerIP string, workerPort string, serverIP string, serverPort string, joinKey string) bool {
	fmt.Println("Joining the cluster...")

	endpoint := "node/join"
	reqJoin := NodeJoinReq{NodeIP: workerIP, NodePort: workerPort, JoinKey: joinKey}
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(reqJoin)

	resp, err := reqServer(endpoint, reqBody)
	if err != nil {
		panic(err)
	}

	respBody := Resp{}
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		panic(err)
	}

	if respBody.Success {
		fmt.Println(respBody.Result)
		fmt.Println("Joined the cluster")
		return true
	} else {
		fmt.Println(respBody.Error)
	}

	return false
}
