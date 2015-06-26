package docker

func Stop(containerName string) int {
    return stopContainer(containerName)
}
