package task

import (
	"time"

	"github.com/docker/docker/client"
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
	Indegree      int
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
