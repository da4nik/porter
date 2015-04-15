package consul

import (
    "encoding/json"
    "flag"
    "log"
)

const pathConsulKV = "/kv/"

type kvItem struct {
    CreateIndex int
    ModifyIndex int
    LockIndex   int
    Key         string
    Flags       int
    Value       []byte
}

type ServiceConfig struct {
    Volumes []string
    Ports   []string
    Env     []string
}

func GetServiceConfig(serviceName string) *ServiceConfig {
    body := apiCall(pathConsulKV + getServiceConfigKey(serviceName))

    var parsedJson []kvItem
    if err := json.Unmarshal(body, &parsedJson); err != nil {
        log.Fatal("Unable to parse JSON. ", string(body))
    }

    var serviceConfig ServiceConfig
    if err := json.Unmarshal(parsedJson[0].Value, &serviceConfig); err != nil {
        log.Fatal("Unable to parse service config. ", err)
    }

    return &serviceConfig
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
