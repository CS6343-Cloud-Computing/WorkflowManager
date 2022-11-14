package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type Server struct {
	ServerIP   string
	ServerPort string
}

// Resp - Generic response
type Resp struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

var serverConfig = Server{"192.168.1.92", "8080"}

func reqServer(endpoint string, reqBody io.Reader) (result []byte, err error) {
	url := "http://" + serverConfig.ServerIP + ":" + serverConfig.ServerPort + "/" + endpoint

	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return result, err
	}
	req.Header.Set("Connection", "close")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	reader := bufio.NewReader(resp.Body)
	resBody, _ := ioutil.ReadAll(reader)
	resp.Body.Close()
	if err == nil {
		return resBody, err
	}

	return nil, errors.New("error handling the error, all hope is lost")
}

func checkServerStatus(cCtx *cli.Context) error {
	endpoint := "server/status"

	resp, err := reqServer(endpoint, nil)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	respBody := Resp{}
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		panic(err)
	}

	if respBody.Success {
		fmt.Println(respBody.Result)
		return nil
	} else {
		fmt.Println(respBody.Error)
	}

	return errors.New("error handling the error, all hope is lost")
}

func submitWorkflow(cCtx *cli.Context) error {
	workflowFile := cCtx.Args().First()

	if workflowFile == "" {
		workflowFile = "./config.yml"
	}

	f, err := os.Open(workflowFile)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	configReader := bufio.NewReader(f)

	endpoint := "workflow/submit"

	resp, err := reqServer(endpoint, configReader)
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
		return nil
	} else {
		fmt.Println(respBody.Error)
	}

	return errors.New("error handling the error, all hope is lost")
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "workflow",
				Usage: "Commands for managing workflow",
				Subcommands: []*cli.Command{
					{
						Name:   "submit",
						Usage:  "Submit workflow for execution",
						Action: submitWorkflow,
					},
					{
						Name:  "ls",
						Usage: "List all the workflows in a cluster.",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("submitted workflow: ", cCtx.Args().First())
							return nil
						},
					},
				},
			},
			{
				Name:  "server",
				Usage: "Commands for managing the Workflow Manager Server",
				Subcommands: []*cli.Command{
					{
						Name:   "status",
						Usage:  "Check the status of the workflow manager",
						Action: checkServerStatus,
					},
				},
			},
			{
				Name:  "worker",
				Usage: "Commands to list all workers",
				Subcommands: []*cli.Command{
					{
						Name:  "ls",
						Usage: "List all the workers in the cluster",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("worker: ", cCtx.Args().First())
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
