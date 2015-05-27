package consul

import (
    "errors"
    "fmt"
    "github.com/hashicorp/consul/api"
    "log"
)

func getServiceConfigKey(service string) string {
    return fmt.Sprintf("service/%s/config", service)
}

func getClient() (client *api.Client, err error) {
    config := api.DefaultConfig()
    client, err = api.NewClient(config)
    if err != nil {
        log.Fatal("Unable to connect to consul: ", err)
        return
    }
    return
}

func getServices(serviceName string) (services []*api.CatalogService, err error) {
    client, err := getClient()
    if err != nil {
        return
    }
    catalog := client.Catalog()
    services, _, err = catalog.Service(serviceName, "", nil)
    if err != nil {
        log.Fatal("Can't get services: ", err)
        return
    }
    return
}

func getServiceConfig(serviceName string) (value []uint8, err error) {
    client, err := getClient()
    if err != nil {
        return
    }
    kv := client.KV()
    pair, _, err := kv.Get(getServiceConfigKey(serviceName), nil)
    if err != nil {
        log.Fatal("Can't get service config ", err)
        return
    }
    if pair == nil {
        err = errors.New(fmt.Sprint("No config for service ", serviceName))
        return
    }
    value = pair.Value
    return
}
