package consul

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/da4nik/porter/docker"
    "github.com/docker/docker/nat"
    consulapi "github.com/hashicorp/consul/api"
    "net/url"
)

const (
    configPrefix = "service"
)

type ContainerConfig struct {
    Volumes []string
    Ports   []string
    Env     []string
}

type ServiceConfig struct {
    ContainerConfig

    Name        string
    LastCommit  string
    Repo        string
    Token       string
    ModifyIndex uint64
    consulData  *consulapi.KVPair
}

func (s *ServiceConfig) serialize() (result []byte, err error) {
    return json.Marshal(s)
}

func (s *ServiceConfig) deserialize(data []byte) error {
    return json.Unmarshal(data, s)
}

func (s *ServiceConfig) Key() string {
    return fmt.Sprintf("%s/%s/config", configPrefix, s.Name)
}

func (s *ServiceConfig) GetModifyIndex() uint64 {
    return s.ModifyIndex
}

func (s *ServiceConfig) SetModifyIndex(index uint64) {
    s.ModifyIndex = index
}

func (s *ServiceConfig) CloneUrl() (cu string, err error) {
    u, err := url.Parse(s.Repo)
    if err != nil {
        return
    }
    u.User = url.User(s.Token)
    if len(s.LastCommit) > 0 {
        u.Fragment = s.LastCommit
    }
    cu = u.String()
    return
}

func (s *ServiceConfig) ContainerName() string {
    return s.Name
}

func (s *ServiceConfig) ImageName() string {
    return s.Name
}

func (s *ServiceConfig) getPorts() (bindings map[nat.Port][]nat.PortBinding, err error) {
    _, bindings, err = nat.ParsePortSpecs(s.Ports)
    return
}

func (s *ServiceConfig) Update() error {
    return SaveConfig(s)
}

func (s *ServiceConfig) Register() error {
    var err error
    if !docker.ContainerIsRunning(s.ContainerName()) {
        return errors.New(fmt.Sprintf("Container '%s' isn't running\n", s.ContainerName()))
    }
    portMap := docker.ContainerPorts(s.Name)
    for intPort, hostPort := range portMap {
        serviceId := fmt.Sprintf("%s:%s", s.Name, intPort)
        err = api.registerService(s.Name, serviceId, hostPort)
        if err != nil {
            s.Deregister()
            return err
        }
    }
    return nil
}

func (s *ServiceConfig) Deregister() error {
    services, err := api.listAgentServices()
    if err != nil {
        return err
    }
    for id, service := range services {
        if service.Service == s.Name {
            err = api.deregisterService(id)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func GetServiceConfig(serviceName string) (serviceConfig *ServiceConfig, err error) {
    serviceConfig = new(ServiceConfig)
    serviceConfig.Name = serviceName
    err = LoadConfig(serviceConfig)
    return
}

func UpdateServiceConfig(serviceName string, repo *url.URL, token *string) error {
    config, err := GetServiceConfig(serviceName)
    if err != nil {
        if _, ok := err.(NoConfigError); !ok {
            return err
        }
        config = new(ServiceConfig)
        config.Name = serviceName
    }
    if repo != nil {
        config.Repo = repo.String()
    }
    if token != nil {
        config.Token = *token
    }
    return config.Update()
}

func ListConfigs() (configList []*ServiceConfig, err error) {
    var pairs consulapi.KVPairs
    pairs, _, err = api.Client().KV().List(configPrefix, nil)
    if err != nil {
        return
    }
    for _, pair := range pairs {
        config := new(ServiceConfig)
        err = fillFromKV(config, pair)
        if err != nil {
            return
        }
        configList = append(configList, config)
    }
    return
}
