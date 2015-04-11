package docker

import (
    "flag"
    "log"
)


func Run() {
    if flag.Arg(1) == "" {
        log.Fatal("No service name given")
    }

    serviceName := flag.Arg(1)

    statusCode := createContainer(serviceName)
    log.Println(statusCode)
    statusCode = startContainer(serviceName)
    log.Println(statusCode)
    inspectContainer(serviceName)
}
