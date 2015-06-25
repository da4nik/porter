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

type NoConfigError struct {
    name, key string
}

func (e NoConfigError) Error() string {
    return fmt.Sprintf("No config for service %s (key %s)\n", e.name, e.key)
}

type ContainerConfig struct {
    Volumes []string
    Ports   []string
    Env     []string
}

type ServiceConfig struct {
    ContainerConfig

    Name       string
    LastCommit string
    Repo       string
    Token      string
    consulData *consulapi.KVPair
}

func (s *ServiceConfig) Key() string {
    return fmt.Sprintf("service/%s/config", s.Name)
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

func (s *ServiceConfig) Update() error {
    s.consulData = new(consulapi.KVPair)
    pair := s.consulData
    pair.Key = s.Key()
    value, err := json.Marshal(s)
    if err != nil {
        return err
    }
    pair.Value = value
    return api.PutKVPair(pair)
}

func (s *ServiceConfig) Delete() error {
    return api.DeleteKvPair(s.Key())
}

func GetServiceConfig(serviceName string) (serviceConfig *ServiceConfig, err error) {
    serviceConfig = new(ServiceConfig)
    serviceConfig.Name = serviceName
    pair, err := api.GetKVPair(serviceConfig.Key())
    if err != nil {
        return
    }
    if pair == nil {
        err = NoConfigError{serviceName, serviceConfig.Key()}
        return
    }
    serviceConfig.consulData = pair
    if err = json.Unmarshal(serviceConfig.consulData.Value, &(serviceConfig)); err != nil {
        return
    }

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

func Config(serviceName string) {
    serviceConfig, err := GetServiceConfig(serviceName)
    if err != nil {
        logger.Fatal(err)
    }
    logger.Println(serviceConfig)
    logger.Println("Volumes: ", serviceConfig.Volumes)
    logger.Println("Ports: ", serviceConfig.Ports)
    logger.Println("Env: ", serviceConfig.Env)
}
