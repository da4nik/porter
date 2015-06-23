package consul

import (
    "github.com/da4nik/porter/docker"
    consulapi "github.com/hashicorp/consul/api"
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
    docker.Build(c.ImageName(), cloneUrl)
    restartEvent := Event{
        consulapi.UserEvent{
            Name: RestartContainersEventName,
        },
    }
    if err != nil {
        return err
    }
    restartEvent.SetServiceName(c.Name)
    restartEvent.ServiceFilter = h.event.ServiceFilter
    return restartEvent.Fire()
}
