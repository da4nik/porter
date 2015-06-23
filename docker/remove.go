package docker

func Remove(containerName string) {
    logger.Printf("Stop container %s (code: %d)\n", containerName, stopContainer(containerName))
    logger.Printf("Remove container %s (code: %d)\n", containerName, removeContainer(containerName))
}
