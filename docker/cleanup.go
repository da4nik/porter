package docker

import (
  "log"
)

const (
  pathUntaggedImages = "/images/json?filters={\"dangling\":[\"true\"]}"
  pathExitedContainers = "/containers/json?all=1&filters={\"status\":[\"exited\"]}"

  pathRemoveContainers = "/containers/" // DELETE
  pathRemoveImages = "/images/" // DELETE
)

func removeEntity(path, id string) {
  deleteApiCall(path + id + "?force=1")
}

func Cleanup() {
  var json []interface{}

  json = getApiCall(pathExitedContainers)
  log.Printf("Found %d exited containers.\n", len(json))

  for i := 0; len(json) > 0 && i < len(json); i++ {
    container := json[i].(map[string]interface{})
    removeEntity(pathRemoveContainers, container["Id"].(string))
  }

  json = getApiCall(pathUntaggedImages)
  log.Printf("Found %d untagged images.\n", len(json))

  for i := 0; len(json) > 0 && i < len(json); i++ {
    image := json[i].(map[string]interface{})
    removeEntity(pathRemoveImages, image["Id"].(string))
  }
}
