package docker

func Push(serviceName, tag string) {
    var statusCode int
    statusCode = tagImage("localhost:5000", serviceName, tag)
    logger.Println(statusCode)
    statusCode = pushImage("localhost:5000", serviceName, tag)
    logger.Println(statusCode)
}
