package task

import (
	"github.com/google/uuid"
	"github.com/docker/go-connections/nat"
	"time"
)

type State int

const (
	Pending State = iota
	Scheduled
	Completed
	Running
	Failed
)

type Task struct {
	ID	uuid.UUID
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
}

type TaskEvent struct {
	ID	uuid.UUID
	State	State
	Timestamp time.Time
	Task Task
}

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
	ctx := context.Backgroud()
	reader,err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("Error pulling the image %s: %v \n", d.Config.Image, err)
		return DockerResult(Error: err)
	}
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
			Name: d.Config.RestartPolicy,
	}
	
	r := container.Resources{
			Memory: d.Config.Memory,
	}

	cc := container.Config{
			Image: d.Config.Image,
			Env: d.Config.Env,
	}

	hc := container.HostConfig{
			RestartPolicy: rp,
			Resources: r,
			PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v \n", d.Config.Image, err)
		return DockerResult{Error:err}
	}

	err := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error starting container using image %s: %v \n", resp.ID, err)
		return DockerResult{Error:err}
	}

	d.ContainerId = resp.ID
	out, err := d.Client.ContainerLogs(ctx, resp.ID,
	                              types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err != null {
	    log.Printf("Error getting logs for container %s: %v\n", resp.ID
	    return DockerResult{Error: err})
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
	    ContainerId: resp.ID,
	    Action: "start",
	    Result: "success"
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

