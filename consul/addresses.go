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
    ca := NewApi()
    service, err := ca.GetService(serviceName)
    if err != nil {
        log.Fatal("Error get service: ", err)
    }
    fmt.Printf("%s:%d\n", service.Address, service.ServicePort)
}
