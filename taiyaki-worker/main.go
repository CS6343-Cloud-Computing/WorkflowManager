package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type NodeJoinReq struct {
	NodeIP   string
	NodePort string
	JoinKey  string
}

var workerIP string
var workerPort string
var serverIP string
var serverPort string
var joinKey string

func reqServer(endpoint string, reqBody io.Reader) (resBody []byte, err error) {
	c := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: c}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "close")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(resp.Body)
	resBody, _ = ioutil.ReadAll(reader)
	resp.Body.Close()

	return resBody, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	if workerIP == "" {
		if env := os.Getenv("WORKER_IP"); env == "" {
			conn, err := net.Dial("udp", "8.8.8.8:80")
			if err != nil {
				log.Fatal(err)
			}
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			workerIP = localAddr.IP.String()
			fmt.Println("No Worker IP provided. But I can do better:", workerIP)
			conn.Close()
		} else {
			workerIP = env
		}
	}

	if workerPort == "" {
		if env := os.Getenv("WORKER_PORT"); env == "" {
			log.Fatal("No server to connect to... :/")
		} else {
			workerPort = env
		}
	}

	if serverIP == "" {
		if env := os.Getenv("SERVER_IP"); env == "" {
			log.Fatal("No server to connect to... :/")
		} else {
			serverIP = env
		}
	}

	if serverPort == "" {
		if env := os.Getenv("SERVER_PORT"); env == "" {
			log.Fatal("No server to connect to... :/")
		} else {
			serverPort = env
		}
	}

	if joinKey == "" {
		if env := os.Getenv("JOIN_KEY"); env == "" {
			log.Fatal("Please provide join key for this cluster")
		} else {
			joinKey = env
		}
	}

	NodeJoin(workerIP, workerPort, serverIP, serverPort, joinKey)
}
