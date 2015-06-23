package docker

import (
    "encoding/json"
    "fmt"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/nat"
    "github.com/docker/docker/runconfig"
    "io/ioutil"
)

const (
    pathCreateContainer  = "/containers/create"    // POST
    pathInspectContainer = "/containers/%s/json"   // GET
    pathStartContainer   = "/containers/%s/start"  // POST
    pathStopContainer    = "/containers/%s/stop"   // POST
    pathRenameContainer  = "/containers/%s/rename" // POST
    pathDeleteContainer  = "/containers/%s"        // DELETE

)

func createContainer(containerName, imageName, tag string, env, volumes, ports []string) int {

    _, bindings, err := nat.ParsePortSpecs(ports)
    if err != nil {
        logger.Println(err)
        return 0
    }

    config := &runconfig.Config{
        Image: imageName,
        Env:   env,
    }
    hostConfig := &runconfig.HostConfig{
        Binds:        volumes,
        PortBindings: bindings,
    }
    mergedConfig := runconfig.MergeConfigs(config, hostConfig)

    params_json, err := json.Marshal(mergedConfig)
    if err != nil {
        logger.Fatal(err)
    }

    resp, err := apiCall("POST", pathCreateContainer, query{"name": {containerName}}, params_json)

    if err != nil {
        logger.Fatal("Unable to create container. ", err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Fatal(err)
    }
    _ = body
    return resp.StatusCode
}

func startContainer(containerName string) int {
    resp, err := apiCall("POST", fmt.Sprintf(pathStartContainer, containerName), nil, nil)
    if err != nil {
        logger.Fatal("Unable to start container. ", err)
    }
    return resp.StatusCode
}

func renameContainer(containerName, newName string) int {
    resp, err := apiCall("POST", fmt.Sprintf(pathRenameContainer, containerName), query{"name": {newName}}, nil)
    if err != nil {
        logger.Fatal("Unable to rename container ", containerName, ". ", err)
    }
    return resp.StatusCode
}

func stopContainer(containerName string) int {
    resp, err := apiCall("POST", fmt.Sprintf(pathStopContainer, containerName), nil, nil)
    if err != nil {
        logger.Fatal("Unable to stop container ", containerName, ". ", err)
    }
    return resp.StatusCode
}

func removeContainer(containerName string) int {
    resp, err := apiCall("DELETE", fmt.Sprintf(pathDeleteContainer, containerName), nil, nil)
    if err != nil {
        logger.Fatal("Unable to delete container ", containerName, ". ", err)
    }
    return resp.StatusCode
}

func inspectContainer(containerName string) types.ContainerJSON {
    var result types.ContainerJSON
    resp, err := apiCall("GET", fmt.Sprintf(pathInspectContainer, containerName), nil, nil)
    if err != nil {
        logger.Fatal("Unable to inspect container ", containerName, ". ", err)
    }
    defer resp.Body.Close()
    r := json.NewDecoder(resp.Body)
    err = r.Decode(&result)
    if err != nil {
        logger.Fatal(err)
    }
    return result
}

func ContainerIsRunning(containerName string) bool {
    return inspectContainer(containerName).State.Running
}
