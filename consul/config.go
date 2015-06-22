package consul

import (
    "encoding/json"
    "fmt"
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
    Name       string
    LastCommit string
    Repo       string
    Token      string

    ContainerConfig
    *consulapi.KVPair
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
    cu = u.String()
    return
}

func (s *ServiceConfig) ContainerName() string {
    // TODO: check for valid symbols
    if s.LastCommit == "" {
        return s.Name
    }
    return fmt.Sprintf("%s.%s", s.Name, s.LastCommit)
}

func (s *ServiceConfig) ImageName() string {
    return s.Name
}

func (s *ServiceConfig) getPorts() (map[nat.Port]struct{}, error) {
    result := make(map[nat.Port]struct{})
    ports, _, err := nat.ParsePortSpecs(s.Ports)
    if err != nil {
        return result, err
    }
    for port, data := range ports {
        result[port] = data
    }
    return result, err
}

func (s *ServiceConfig) SeriveIds() (map[nat.Port]string, error) {
    result := make(map[nat.Port]string)
    ports, err := s.getPorts()
    if err != nil {
        return result, err
    }
    for port, _ := range ports {
        result[port] = fmt.Sprintf("%s:%s", s.Name, port)
    }
    return result, err
}

func (s *ServiceConfig) Register() error {
    ids, err := s.SeriveIds()
    if err != nil {
        return err
    }
    for port, id := range ids {
        err = api.registerService(s, id, port.Int())
        if err != nil {
            s.Deregister()
            return err
        }
    }
    return nil
}

func (s *ServiceConfig) Deregister() error {
    ids, err := s.SeriveIds()
    if err != nil {
        return err
    }
    for _, id := range ids {
        err = api.deregisterService(id)
        if err != nil {
            s.Deregister()
            return err
        }
    }
    return nil
}

func (s *ServiceConfig) Update() error {
    s.KVPair = new(consulapi.KVPair)
    pair := s.KVPair
    pair.Key = s.Key()
    value, err := json.Marshal(s)
    if err != nil {
        return err
    }
    pair.Value = value
    return api.PutKVPair(pair)
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
    serviceConfig.KVPair = pair
    if err = json.Unmarshal(serviceConfig.Value, &(serviceConfig)); err != nil {
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
