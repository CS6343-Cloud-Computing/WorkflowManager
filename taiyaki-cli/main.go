package main

import (
	"bufio"
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

//Resp - Generic response
type Resp struct {
	Result  string
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

var serverConfig = Server{"192.168.1.92", "8080"}

func reqServer(endpoint string, reqBody io.Reader) (result string, err error) {
	url := "http://" + serverConfig.ServerIP + ":" + serverConfig.ServerPort + "/" + endpoint

	req, err := http.NewRequest(http.MethodGet, url, reqBody)
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
		return string(resBody), err
	}

	return "", errors.New("error handling the error, all hope is lost")
}

func checkServerStatus(cCtx *cli.Context) error {
	endpoint := "server/status"
	resp, err := reqServer(endpoint, nil)
	fmt.Println(resp)
	return err
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "workflow",
				Usage: "Commands for managing workflow",
				Subcommands: []*cli.Command{
					{
						Name:  "submit",
						Usage: "Submit workflow for execution",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("submitted workflow: ", cCtx.Args().First())
							return nil
						},
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
