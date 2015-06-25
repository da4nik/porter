package consul

import (
    "github.com/da4nik/porter/docker"
)

const (
    BuildImageEventName = "build-image"
)

type BuildImageHandler struct {
    event *Event
}

func (h *BuildImageHandler) Handle() error {
    c, err := GetServiceConfig(h.event.ServiceName())
    if err != nil {
        return err
    }
    cloneUrl, err := c.CloneUrl()
    if err != nil {
        return err
    }
    err = docker.Build(c.ImageName(), cloneUrl)
    if err != nil {
        return err
    }
    restartEvent := NewRestartContainerEvent(c.Name, h.event.NodeFilter, h.event.ServiceFilter)
    return restartEvent.Fire()
}

func NewBuildImageEvent(serviceName, nodeFilter, serviceFilter string) Event {
    return newEvent(BuildImageEventName, serviceName, nodeFilter, serviceFilter)
}
