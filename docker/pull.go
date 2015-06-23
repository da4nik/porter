package docker

func Pull(serviceName, tag string) {
    var statusCode int
    statusCode = createImage("localhost:5000", serviceName, tag)
    logger.Println(statusCode)
}
