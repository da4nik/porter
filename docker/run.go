package docker

import (
  "log"
  "github.com/da4nik/porter/consul"
)

func RunContainer(serviceName string) {
  serviceConfig := consul.GetServiceConfig(serviceName)

  log.Println(string(serviceConfig.Value))
}
