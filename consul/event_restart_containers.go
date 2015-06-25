package consul

import (
    "fmt"
    "github.com/da4nik/porter/docker"
    "time"
)

const (
    RestartContainersEventName = "restart-containers"
)

type StartContainerError string

func (e StartContainerError) Error() string {
    return fmt.Sprintf("Error start container: %s\n", string(e))
}

type RestartContainersHandler struct {
    event *Event
}

func (h *RestartContainersHandler) Handle() error {
    c, err := GetServiceConfig(h.event.ServiceName())
    if err != nil {
        return err
    }
    old_name := c.ContainerName() + "_old"
    c.Deregister()
    docker.Rename(c.ContainerName(), old_name)
    docker.Stop(old_name)
    docker.Run(c.ContainerName(), c.ImageName(), c.LastCommit, c.Env, c.Volumes, c.Ports)
    <-time.After(10 * time.Second)
    if !docker.ContainerIsRunning(c.ContainerName()) {
        docker.Remove(c.ContainerName())
        docker.Rename(old_name, c.ContainerName())
        docker.Start(old_name)
        return StartContainerError("New container isn't started normally")
    }
    c.Register()
    docker.Remove(old_name)
    return nil
}

func NewRestartContainerEvent(serviceName, nodeFilter, serviceFilter string) Event {
    return newEvent(RestartContainersEventName, serviceName, nodeFilter, serviceFilter)
}
