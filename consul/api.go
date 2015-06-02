package consul

import (
    "errors"
    "fmt"
    "github.com/hashicorp/consul/api"
    "log"
)

type ConsulApi struct {
    config *api.Config
    client *api.Client
}

func NewApi() *ConsulApi {
    var (
        err error
    )
    c := new(ConsulApi)
    c.config = api.DefaultConfig()
    c.client, err = api.NewClient(c.config)
    if err != nil {
        log.Fatal("Can't init consul client: ", err)
    }
    return c
}

func (c *ConsulApi) GetService(serviceName string) (service *api.CatalogService, err error) {
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

func (c *ConsulApi) GetServiceConfig(serviceName string) (value []uint8, err error) {
    kv := c.client.KV()
    pair, _, err := kv.Get(getServiceConfigKey(serviceName), nil)
    if err != nil {
        return
    }
    if pair == nil {
        err = errors.New(fmt.Sprint("No config for service ", serviceName))
        return
    }
    value = pair.Value
    return
}

func getServiceConfigKey(service string) string {
    return fmt.Sprintf("service/%s/config", service)
}
