package consul

import (
    "flag"
    "fmt"
    "log"
)

func Addresses() {
    if flag.Arg(1) == "" {
        log.Fatal("No service name given")
    }

    serviceName := flag.Arg(1)
    services, err := getServices(serviceName)
    if err != nil {
        return
    }
    for _, service := range services {
        fmt.Printf("%s:%d\n", service.Address, service.ServicePort)
    }
}
