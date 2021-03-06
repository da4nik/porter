package docker

import (
    "sync"
)

const (
    pathUntaggedImages   = `/images/json?filters={"dangling":["true"]}`
    pathExitedContainers = `/containers/json?all=1&filters={"status":["exited"]}`

    pathRemoveContainers = "/containers/" // DELETE
    pathRemoveImages     = "/images/"     // DELETE
)

func removeEntity(path, id string, wg *sync.WaitGroup) {
    defer wg.Done()
    apiCall("DELETE", path+id, query{"force": {"1"}}, nil)
}

func removeContainersOrImages(queryPath, deletePath, messageTpl string) {
    var json []interface{}
    var wg sync.WaitGroup

    json = jsonGetApiCall(queryPath)
    logger.Printf(messageTpl, len(json))

    for i := 0; len(json) > 0 && i < len(json); i++ {
        container := json[i].(map[string]interface{})
        wg.Add(1)
        go removeEntity(deletePath, container["Id"].(string), &wg)
    }
    wg.Wait()
}

func Cleanup() {
    removeContainersOrImages(pathExitedContainers, pathRemoveContainers, "Found %d exited containers.\n")
    removeContainersOrImages(pathUntaggedImages, pathRemoveImages, "Found %d untagged images.\n")
}
