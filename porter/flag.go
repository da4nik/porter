package main

import (
    "flag"
    "fmt"
    "os"
    "github.com/da4nik/porter/utils"
    "strings"
)

var (
    flDebug = *flag.Bool("D", false, "Set debug mode")
    Action  string
    ServiceName string
)

var allActions  = []string{
    "cleanup",
    "run",
    "config",
}

var actionsWithService = []string{
    "run",
    "config",
}

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

func parseParams() {
    flag.Parse()

    if flag.NArg() == 0 {
        fmt.Fprintf(os.Stdout, "No action provided.\n\n")
        usage()
        os.Exit(1)
    }
    Action = strings.ToLower(flag.Arg(0))

    if result := utils.IsArrayIncludes(actionsWithService, Action); result != -1 {
        if flag.NArg() < 2 {
            fmt.Fprintf(os.Stdout, "No service name provided for action.\n\n")
            usage()
            os.Exit(1)
        }
        ServiceName = strings.ToLower(flag.Arg(1))
    }
}
