package docker

func Build(serviceName, remoteUrl string) {
    tag := serviceName
    statusCode := buildImage(tag, remoteUrl)
    logger.Println(statusCode)
}
