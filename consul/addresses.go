package consul

import (
    "fmt"
)

func Addresses(serviceName string) {
    service, err := api.GetService(serviceName)
    if err != nil {
        logger.Fatal("Error get service: ", err)
    }
    fmt.Printf("%s:%d\n", service.Address, service.ServicePort)
}
