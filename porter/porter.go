package main

import (
    "os"

    "github.com/da4nik/porter/consul"
    "github.com/da4nik/porter/docker"
)

func main() {
    parseParams()

    switch Action {
    case "cleanup":
        docker.Cleanup()
    case "run":
        docker.Run()
    case "addresses":
        consul.Addresses()
    case "config":
        consul.Config()
    default:
        usage()
        os.Exit(1)
    }
}
