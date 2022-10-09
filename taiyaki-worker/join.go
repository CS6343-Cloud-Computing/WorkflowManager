package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func NodeJoin(workerIP string, workerPort string, serverIP string, serverPort string, joinKey string) string {
	fmt.Println("Starting Join Process...")

	reqJoin :=  NodeJoinReq{NodeIP: workerIP, NodePort: workerPort, JoinKey: joinKey}
	reqBuffer := new(bytes.Buffer)
	json.NewEncoder(reqBuffer).Encode(reqJoin)
	fmt.Println("Sending initial join request...")

	resp, err := reqServer("http://"+serverIP+":8080/node/join", reqBuffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp))

	// type NodeInitialJoinResp struct {
	// 	CaCert     string
	// 	ClientCert string
	// 	Success    bool   `json:"success"`
	// 	Error      string `json:"error"`
	// }

	// respj := NodeInitialJoinResp{}
	// err = json.Unmarshal(resp, &respj)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }

	// fmt.Println("Received response: ", respj)

	// //Step 3
	// if !respj.Success {
	// 	panic(respj.Error)
	// }
	return ""
}
