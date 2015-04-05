package docker

import (
  "log"
  "flag"

  "github.com/da4nik/porter/consul"
)

func RunContainer() {
  if flag.Arg(1) == "" {
    log.Fatal("No service name given")
  }

  serviceConfig := consul.GetServiceConfig(flag.Arg(1))

  log.Println(string(serviceConfig.Value))
}
