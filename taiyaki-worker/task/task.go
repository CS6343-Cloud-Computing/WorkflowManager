package task

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	Cmd           []string
	Image         string
	Env           []string
	RestartPolicy string
	Query         string
}

type Task struct {
	ID            uuid.UUID
	WorkflowID    string
	ContainerId   string
	Name          string
	State         State
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
	Config        Config
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

type Docker struct {
	Client      *client.Client
	Task        Task
	ContainerId string
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(ctx, d.Task.Config.Image, types.ImagePullOptions{})

	if err != nil {
		log.Printf("Error pulling the image %s: %v \n", d.Task.Config.Image, err)
		return DockerResult{Error: err}
	}

	io.Copy(os.Stdout, reader)
	rp := container.RestartPolicy{
		Name: d.Task.Config.RestartPolicy,
	}

	cc := container.Config{
		Image: d.Task.Config.Image,
		Env:   d.Task.Config.Env,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Task.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v \n", d.Task.Config.Image, err)
		return DockerResult{Error: err}
	}

	err2 := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err2 != nil {
		log.Printf("Error starting container using image %s: %v \n", resp.ID, err2)
		return DockerResult{Error: err}
	}

	d.ContainerId = resp.ID
	out, err := d.Client.ContainerLogs(ctx, resp.ID,
		types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerId: resp.ID,
		Action:      "start",
		Result:      "success",
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
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func ValidStateTransition(src State, dst State) bool {
	return Contains(stateTransitionMap[src], dst)
}

func NewConfig(task *Task) Config {
	config := &task.Config
	return *config
}

func NewDocker(config Config) *Docker {
	d := new(Docker)
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d.Client = dc
	d.Task.Config = config
	return d
}
