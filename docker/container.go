package docker

import (
    "log"
    "fmt"
    "encoding/json"

    "github.com/da4nik/porter/consul"
)

const (
    pathCreateContainer = "/containers/create" // POST
    pathInspectContainer = "/containers/%s/json" // GET
    pathStartContainer = "/containers/%s/start" // POST
)

func createContainer(serviceName string) int {

    serviceConfig := consul.GetServiceConfig(serviceName)

    parameters := map[string]interface{}{
        "Image": serviceName,
        "Env": serviceConfig.Env,
        "HostConfig": map[string]interface{}{
            "Binds": serviceConfig.Volumes,
        },
    }

    params_json, _ := json.Marshal(parameters)

    resp, err := apiCall("POST", pathCreateContainer, []string{"name="+serviceName, }, params_json)

    if err != nil {
        log.Fatal("Unable to create container. ", err)
    }

    return resp.StatusCode
}

func startContainer(containerName string) int {
    resp, err := apiCall("POST", fmt.Sprintf(pathStartContainer, containerName), nil, nil)
    if err != nil {
        log.Fatal("Unable to start container. ", err)
    }
    return resp.StatusCode
}

func inspectContainer(containerName string) interface{} {
    resp, err := apiCall("GET", fmt.Sprintf(pathInspectContainer, containerName), nil, nil)
    if err != nil {
        log.Fatal("Unable to inspect container ", containerName, ". ", err)
    }
    return resp
}

