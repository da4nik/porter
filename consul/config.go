package consul

import (
    "encoding/json"
    "flag"
    "log"
)

type ServiceConfig struct {
    Volumes []string
    Ports   []string
    Env     []string
}

func GetServiceConfig(serviceName string) *ServiceConfig {

    serviceConfig := new(ServiceConfig)
    value, err := getServiceConfig(serviceName)

    if err != nil {
        log.Fatal("Unable to get service config. ", err)
    }
    if err := json.Unmarshal(value, serviceConfig); err != nil {
        log.Fatal("Unable to parse service config. ", err)
    }

    return serviceConfig
}

func Config() {
    if flag.Arg(1) == "" {
        log.Fatal("No service name given")
    }

    serviceName := flag.Arg(1)
    serviceConfig := GetServiceConfig(serviceName)
    log.Println(serviceConfig)
    log.Println("Volumes: ", serviceConfig.Volumes)
    log.Println("Ports: ", serviceConfig.Ports)
    log.Println("Env: ", serviceConfig.Env)
}
