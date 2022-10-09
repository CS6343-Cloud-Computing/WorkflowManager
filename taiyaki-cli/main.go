package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func reqServer(endpoint string, reqBody io.Reader) (resBody []byte, err error) {

	return nil, errors.New("error handling the error, all hope is lost")
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
						Name:  "status",
						Usage: "Check the status of the workflow manager",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("Status of the Workflow Manager: ", cCtx.Args().First())
							return nil
						},
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
