package consul

import (
    "errors"
    "fmt"
    consulapi "github.com/hashicorp/consul/api"
    "strconv"
    "time"
)

type NotUpdatedError struct {
    key string
}

func (e NotUpdatedError) Error() string {
    return fmt.Sprintf("Record '%s' already chaged", e.key)
}

var api *ConsulApi

type ConsulApi struct {
    config *consulapi.Config
    client *consulapi.Client
}

func init() {
    GetApi()
}

func GetApi() *ConsulApi {
    var err error
    if api == nil {
        api = new(ConsulApi)
        api.config = consulapi.DefaultConfig()
        api.client, err = consulapi.NewClient(api.config)
        if err != nil {
            logger.Fatal("Can't init consul client: ", err)
        }

    }
    return api
}

func (c *ConsulApi) Client() *consulapi.Client {
    return c.client
}

func (c *ConsulApi) GetService(serviceName string) (service *consulapi.CatalogService, err error) {
    catalog := c.client.Catalog()
    services, _, err := catalog.Service(serviceName, "", nil)
    if err != nil {
        return
    }
    if len(services) == 0 {
        err = errors.New(fmt.Sprintf("Can't find service '%s'\n", serviceName))
        return
    }
    service = services[0]
    return
}

func (c *ConsulApi) GetKVPair(key string) (pair *consulapi.KVPair, err error) {
    kv := c.client.KV()
    pair, _, err = kv.Get(key, nil)
    return
}

func (c *ConsulApi) PutKVPair(pair *consulapi.KVPair) error {
    var (
        err    error
        kv     *consulapi.KV
        stored bool
    )
    kv = c.client.KV()
    if pair.ModifyIndex != 0 {
        stored, _, err = kv.CAS(pair, nil)
        if err != nil {
            return err
        }
        if !stored {
            return NotUpdatedError{pair.Key}
        }
    } else {
        _, err = kv.Put(pair, nil)
    }
    return err
}

func (c *ConsulApi) DeleteKvPair(key string) error {
    _, err := c.client.KV().Delete(key, nil)
    return err
}

func (c *ConsulApi) registerService(serviceName, serviceID, serviceHostPort string) error {
    registration := new(consulapi.AgentServiceRegistration)
    registration.ID = serviceID
    registration.Name = serviceName
    port, _ := strconv.Atoi(serviceHostPort)
    registration.Port = port
    agent, err := c.client.Agent().Self()
    if err != nil {
        return err
    }
    registration.Address = agent["Member"]["Addr"].(string)
    return c.client.Agent().ServiceRegister(registration)
}

func (c *ConsulApi) deregisterService(serviceID string) error {
    return c.client.Agent().ServiceDeregister(serviceID)
}

func (c *ConsulApi) listAgentServices() (map[string]*consulapi.AgentService, error) {
    return c.Client().Agent().Services()
}

func (c *ConsulApi) FireEvent(params *Event) (result string, err error) {
    event := c.client.Event()
    result, _, err = event.Fire(&params.UserEvent, new(consulapi.WriteOptions))
    if err != nil {
        return
    }
    return
}

func (c *ConsulApi) ListEvents(waitIndex uint64) (events []*consulapi.UserEvent, lastIndex uint64, err error) {
    var (
        meta *consulapi.QueryMeta
    )
    e := c.client.Event()
    q := &consulapi.QueryOptions{
        WaitIndex: waitIndex,
        WaitTime:  15 * time.Second,
    }
    events, meta, err = e.List("", q)
    if err != nil {
        return
    }
    lastIndex = meta.LastIndex
    // only new events
    for i := 0; i < len(events); i++ {
        if e.IDToIndex(events[i].ID) == waitIndex {
            events = events[i+1:]
            break
        }
    }
    return
}
