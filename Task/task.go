package task

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Completed
	Running
	Failed
)

type Config struct{
	Name string
	AttachStdin bool
	AttachStdout bool
	AttachStderr bool
	Cmd []string
	Image string
	Memory int64
	Disk int64
	Env []string
	RestartPolicy string
}

type Task struct {
	ID	uuid.UUID
	ContainerId	string
	Name	string
	State	State
	Image	string
	Memory	int
	Disk	int
	ExposedPorts	nat.PortSet
	PortBindings	map[string]string
	RestartPolicy	string
	StartTime	time.Time
	FinishTime	time.Time
	Config	Config
}

type TaskEvent struct {
	ID	uuid.UUID
	State	State
	Timestamp time.Time
	Task Task
}

type Docker struct{
	Client *client.Client
	Config Config
	ContainerId string
}

type DockerResult struct{
	Error error
	Action string
	ContainerId string
	Result string
}

func (d *Docker) Run() DockerResult{
	ctx := context.Background()
	fmt.Println("reached 12 " )

	reader,err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})

	if err != nil {
		log.Printf("Error pulling the image %s: %v \n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	fmt.Println("reached 14 ")

	io.Copy(os.Stdout, reader)

	fmt.Println("reached 15 ")

	rp := container.RestartPolicy{
			Name: d.Config.RestartPolicy,
	}
	
	fmt.Println("reached 16 ")

	r := container.Resources{
			Memory: d.Config.Memory,
	}

	fmt.Println("reached 17 ")

	cc := container.Config{
			Image: d.Config.Image,
			Env: d.Config.Env,
	}

	fmt.Println("reached 18 ")

	hc := container.HostConfig{
			RestartPolicy: rp,
			Resources: r,
			PublishAllPorts: true,
	}

	fmt.Println("reached 19 ")

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v \n", d.Config.Image, err)
		return DockerResult{Error:err}
	}

	fmt.Println("reached 20 ")

	err2 := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err2 != nil {
		log.Printf("Error starting container using image %s: %v \n", resp.ID, err2)
		return DockerResult{Error:err}
	}

	d.ContainerId = resp.ID
	out, err := d.Client.ContainerLogs(ctx, resp.ID,
	                              types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err != nil {
	    log.Printf("Error getting logs for container %s: %v\n", resp.ID,err)
	    return DockerResult{Error:err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
	    ContainerId: resp.ID,
	    Action: "start",
	    Result: "success",
	}
}

func (d *Docker) Stop() DockerResult {
	ctx := context.Background()
	log.Printf(
		"Attempting to stop container %v", d.ContainerId)
	err := d.Client.ContainerStop(ctx, d.ContainerId, nil)
	if err != nil {
		panic(err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}

	err = d.Client.ContainerRemove(
		ctx,
		d.ContainerId,
		removeOptions,
	)
	if err != nil {
		panic(err)
	}
	return DockerResult{Action: "stop", Result: "success", Error: nil}
}

var stateTransitionMap = map[State][]State{
	Pending: {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:	{Running, Completed, Failed},
	Completed:	{},
	Failed:	{},
}

func Contains(states []State, state State) bool{
	for _,s := range states{
		if s== state{
			return true
		}
	}
	return false
}

func ValidStateTransition(src State, dst State) bool {
	return Contains(stateTransitionMap[src],dst)
}

func NewConfig(task *Task) Config{
	config := &task.Config
	return *config;
}

func NewDocker(config Config) *Docker{
	d := new(Docker)
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d.Client = dc
	d.Config = config
	return d
}


