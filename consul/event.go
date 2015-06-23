package consul

import (
    "encoding/json"
    consulapi "github.com/hashicorp/consul/api"
    "io"
)

type Handler interface {
    Handle() error
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
}
