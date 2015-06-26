package consul

import (
    consulapi "github.com/hashicorp/consul/api"
    "time"
)

type Handler interface {
    Handle() error
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
    case BuildImageEventName:
        h = &BuildImageHandler{
            event: e,
        }
    case RestartContainersEventName:
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

func newEvent(eventName, serviceName, nodeFilter, serviceFilter string) Event {
    event := Event{
        consulapi.UserEvent{
            Name:          eventName,
            ServiceFilter: serviceFilter,
            NodeFilter:    nodeFilter,
        },
    }
    event.SetServiceName(serviceName)
    return event
}

func ListenEvents() {
    const (
        startRetryTime = 5 * time.Second
        maxRetryTime   = 3 * time.Minute
    )
    var (
        events    []*consulapi.UserEvent
        lastIndex uint64
        err       error
        retryTime = startRetryTime
    )
    for {
        events, lastIndex, err = api.ListEvents(lastIndex)
        if err != nil {
            logger.Print(err)
            logger.Println("Retry after", retryTime)
            <-time.After(retryTime)
            retryTime *= 2
            if retryTime > maxRetryTime {
                retryTime = maxRetryTime
            }
            continue
        }
        retryTime = startRetryTime
        for _, ue := range events {
            event := Event{
                UserEvent: *ue,
            }
            logger.Printf("Start event '%s'\n", event.Name)
            err = event.GetHandler().Handle()
            if err != nil {
                logger.Println(err)
            }
        }
    }
}
