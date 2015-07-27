package main

import (
    "log"
    "os"

    "github.com/da4nik/porter/consul"
    "github.com/da4nik/porter/docker"
    kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
    app     = kingpin.New("porter", "Containers via consul handle utility.")
    flDebug = app.Flag("debug", "Enable debug mode.").Bool()

    cleanup = app.Command("cleanup", "Cleanup untagged images and exited containers.")

    run            = app.Command("run", "Run service container.")
    runServiceName = run.Arg("name", "Name of service").Required().String()
    runServiceTag  = run.Arg("tag", "Tag of image. Default value - last commit id").String()

    addresses          = app.Command("addresses", "Prints all addresses for service in cluster")
    addressServiceName = addresses.Arg("name", "Name of service").Required().String()

    config            = app.Command("config", "Print service config")
    configServiceName = config.Arg("name", "Name of service").Required().String()

    setConfig            = app.Command("set-config", "Create/update service config")
    setConfigServiceName = setConfig.Arg("name", "Name of service").Required().String()
    setConfigRepo        = setConfig.Arg("repo", "Github repository clone url").URL()
    setConfigToken       = setConfig.Arg("token", "Github access token").String()

    setConfigVolumes = app.Command("set-volumes", "Set service container volumes")
    setConfigVolume  = setConfigVolumes.Arg("volume", "Container volumes").Strings()

    setConfigPorts = app.Command("set-ports", "Set service container ports")
    setConfigPort  = setConfigPorts.Arg("port", "Container ports").Strings()

    setConfigEnvs = app.Command("set-envs", "Set service container env")
    setConfigEnv  = setConfigEnvs.Arg("env", "Container env").Strings()

    build            = app.Command("build", "Build docker image")
    buildServiceName = build.Arg("name", "Name of service").Required().String()

    push             = app.Command("push", "Push docker image to registry")
    pushServiceName  = push.Arg("name", "Name of service").Required().String()
    pushLastCommitId = push.Arg("lastCommit", "git last commit").Required().String()

    pull             = app.Command("pull", "Pull image")
    pullServiceName  = pull.Arg("name", "Name of service").Required().String()
    pullLastCommitId = pull.Arg("lastCommit", "git last commit").Required().String()

    listen = app.Command("listen", "Listen consul events")
)

func getConfig(serviceName string) *consul.ServiceConfig {
    config, err := consul.GetServiceConfig(serviceName)
    if err != nil {
        log.Fatal(err)
    }
    return config
}

func main() {
    switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case cleanup.FullCommand():
        docker.Cleanup()

    case run.FullCommand():
        c := getConfig(*runServiceName)
        tag := c.LastCommit
        if runServiceTag != nil {
            tag = *runServiceTag
        }
        docker.Run(c.ContainerName(), c.ImageName(), tag, c.Env, c.Volumes, c.Ports)

    case addresses.FullCommand():
        consul.Addresses(*addressServiceName)

    case config.FullCommand():
        serviceConfig, err := consul.GetServiceConfig(*configServiceName)
        if err != nil {
            log.Fatal(err)
        }
        log.Println(serviceConfig)
        log.Println("Volumes: ", serviceConfig.Volumes)
        log.Println("Ports: ", serviceConfig.Ports)
        log.Println("Env: ", serviceConfig.Env)

    case setConfig.FullCommand():
        consul.UpdateServiceConfig(*setConfigServiceName, *setConfigRepo, setConfigToken)

    case build.FullCommand():
        c := getConfig(*buildServiceName)
        cloneUrl, err := c.CloneUrl()
        if err != nil {
            log.Fatal(err)
        }
        docker.Build(*buildServiceName, cloneUrl)

    case push.FullCommand():
        docker.Push(*pushServiceName, *pushLastCommitId)

    case pull.FullCommand():
        docker.Pull(*pullServiceName, *pullLastCommitId)

    case listen.FullCommand():
        consul.ListenEvents()

    default:
        kingpin.Usage()
    }
    log.Println("Done")
}
