package main

import (
    "os"

    "github.com/da4nik/porter/consul"
    "github.com/da4nik/porter/docker"
    "fmt"
)

func main() {
    parseParams()

    fmt.Println("Action", Action)

    switch Action {
    case "cleanup":
        docker.Cleanup()
    case "run":
        docker.Run()
    case "config":
        consul.Config()
    default:
        usage()
        os.Exit(1)
    }
}
