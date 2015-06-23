package docker

func Run(containerName, imageName, tag string, envs, volumes, ports []string) {

    statusCode := createContainer(containerName, imageName, tag, envs, volumes, ports)
    logger.Printf("Create container '%s' (code: %d)\n", containerName, statusCode)
    statusCode = startContainer(containerName)
    logger.Printf("Start container '%s' (code: %d)\n", containerName, statusCode)
    inspectContainer(containerName)
}
