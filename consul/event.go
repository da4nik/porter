package consul

import (
    "encoding/json"
    "github.com/da4nik/porter/docker"
    consulapi "github.com/hashicorp/consul/api"
    "io"
)

const (
    BuildImage        = "build-image"
    RestartContainers = "restart-containers"
)

type Handler interface {
    Handle() error
}

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
            Name: RestartContainers,
        },
    }
    if err != nil {
        return err
    }
    restartEvent.SetServiceName(c.Name)
    restartEvent.ServiceFilter = h.event.ServiceFilter
    return restartEvent.Fire()
}

type RestartContainersHandler struct {
    event *Event
}

func (h *RestartContainersHandler) Handle() error {
    c, err := GetServiceConfig(h.event.ServiceName())
    if err != nil {
        return err
    }
    old_name := c.Name + "_old"
    c.Deregister()
    docker.Rename(c.Name, old_name)
    docker.Run(c.ContainerName(), c.ImageName(), c.LastCommit, c.Env, c.Volumes, c.Ports)
    c.Register()
    docker.Remove(old_name)
    return nil
}

type NoneHandler struct {
    event *Event
}

func (h *NoneHandler) Handle() error {
    logger.Printf("Unknown event type '%s' with service '%s'\n", h.event.Name, h.event.ServiceName())
    return nil
}

type Event struct {
    consulapi.UserEvent
}

func (e *Event) Fire() error {
    _, err := api.FireEvent(e)
    return err
}

func (e *Event) ServiceName() string {
    return string(e.Payload)
}

func (e *Event) SetServiceName(name string) {
    e.Payload = []byte(name)
}

func (e *Event) GetHandler() (h Handler) {
    switch e.Name {
    case BuildImage:
        h = &BuildImageHandler{
            event: e,
        }
    case RestartContainers:
        h = &RestartContainersHandler{
            event: e,
        }
    default:
        h = &NoneHandler{
            event: e,
        }
    }
    return
}

func ProcessEvents(r io.Reader) {
    var (
        list []Event
    )
    d := json.NewDecoder(r)
    err := d.Decode(&list)
    if err != nil {
        logger.Println(err)
        return
    }
    for _, e := range list {
        logger.Printf("Process event '%s' (payload '%s')\n", e.Name, e.Payload)
        err = e.GetHandler().Handle()
        if err != nil {
            logger.Println(err)
        }
    }
    logger.Println("Done")
}