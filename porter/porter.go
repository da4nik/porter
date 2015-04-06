package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/da4nik/porter/consul"
    "github.com/da4nik/porter/docker"
)

func usage() {
    var help string

    fmt.Fprint(os.Stdout, "Usage: porter [OPTIONS] COMMAND [arg...]\n\nContainers via consul handle utility \n\nOptions:\n")

    commands := [][]string{
        {"cleanup", "Cleanup untagged images and exited containers."},
        {"run", "Run service container"},
    }

    for _, command := range commands {
        help += fmt.Sprintf("    %-10.10s%s\n", command[0], command[1])
    }
    fmt.Fprintf(os.Stdout, "%s\n", help)
}

func main() {

    flag.Parse()

    if len(flag.Args()) == 0 {
        usage()
        os.Exit(1)
    }

    switch strings.ToLower(os.Args[1]) {
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
