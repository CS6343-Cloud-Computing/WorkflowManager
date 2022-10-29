package client

import (
	"bufio"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
)
type Resp struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func ReqServer(serverIP string, serverPort string, endpoint string, reqBody io.Reader) (resBody []byte, err error) {
	url := "http://" + serverIP + ":" + serverPort + "/" + endpoint
	c := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: c}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, reqBody)
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
