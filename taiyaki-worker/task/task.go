package task

import (
	"context"
	"fmt"
	"io"
	"log"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
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
	Indegree      int
	Persistence   bool
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

func (d *Docker) modifyDockerFile(dockerBuildCtxDir string, image string) {
	filePathNew, _ := homedir.Expand(dockerBuildCtxDir)
	input, err := os.ReadFile(filePathNew + "/baseDockerfile")

	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	lines[1] = strings.Replace(lines[1], "$1", image, 1)
	fmt.Println(lines[1])

	lines[2] = strings.Replace(lines[2], "$2", "\""+d.Task.Config.Cmd[0]+"\"", 1)
	fmt.Println(lines[2])

	// lines[3] = strings.Replace(lines[3], "$3", d.Task.ContainerId, 1)
	fmt.Println(lines[3])

	lines[4] = strings.Replace(lines[4], "$4", strconv.Itoa(d.Task.Indegree), 1)
	fmt.Println(lines[4])

	output := strings.Join(lines, "\n")

	err = os.Remove(filePathNew + "/Dockerfile")
	if os.IsNotExist(err) {
		log.Println(err)
	}
	err = os.WriteFile(filePathNew+"/Dockerfile", []byte(output), 0644)

	if err != nil {
		log.Println(err)
	}

}

func (d *Docker) BuildImage(dockerBuildCtxDir, tagName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	dockerFileTarReader := GetContext(dockerBuildCtxDir)
	fmt.Println(dockerFileTarReader)

	resp, err := d.Client.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{tagName},
			NoCache:    false,
			Remove:     false,
		}) //cli is the docker client instance created from the engine-api
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println(err, " :unable to read image build response")
	}
	fmt.Println("BUILDING SUCCESSFULL:", os.Stdout, resp)
	if err != nil {
		log.Println(err, " :unable to build docker image")
		return err
	}
	return nil
}

func GetContext(filePath string) io.Reader {
	// Use homedir.Expand to resolve paths like '~/repos/myrepo'
	filePathNew, _ := homedir.Expand(filePath)
	ctx, _ := archive.TarWithOptions(filePathNew, &archive.TarOptions{})
	return ctx
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	images, err := d.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	isImagePresent := false
	for _, image := range images {
		if image.RepoTags == nil {
			continue
		}
		fmt.Println(image.RepoTags)
		imageName := strings.Split(image.RepoTags[0], ":")
		if imageName == nil {
			imageName = strings.Split(image.RepoTags[0], " ")
		}
		if imageName[0] == "mod_"+d.Task.Config.Image {
			fmt.Println("==============image name is same=============================================================")
			isImagePresent = true
			break
		}
	}
	if !isImagePresent {
		d.modifyDockerFile("~/WorkflowManager/taiyaki-worker/docker/", d.Task.Config.Image)
		d.BuildImage("~/WorkflowManager/taiyaki-worker/docker", "mod_"+d.Task.Config.Image)
	}
	d.Task.Config.Image = "mod_" + d.Task.Config.Image
	//imagepush

	rp := container.RestartPolicy{
		Name: d.Task.Config.RestartPolicy,
	}

	d.Task.Config.Env = []string{"consume_topic=" + d.Task.ContainerId}

	cc := container.Config{
		Image: d.Task.Config.Image,
		Env:   d.Task.Config.Env,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Task.ContainerId)
	if err != nil {
		log.Printf("Error creating container using image %s: %v \n", d.Task.Config.Image, err)
		return DockerResult{Error: err}
	}

	err2 := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err2 != nil {
		log.Printf("Error starting container using image %s: %v \n", resp.ID, err2)
		return DockerResult{Error: err}
	}

	fmt.Println("-------------------------container created with id " + resp.ID)
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
		log.Println("Container already stopped")
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
		log.Println("Container already removed")
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

func NewDocker(config Config, t Task) *Docker {
	d := new(Docker)
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d.Client = dc
	d.Task.Config = config
	d.Task = t
	return d
}
