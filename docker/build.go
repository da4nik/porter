package docker

import (
    "errors"
    "fmt"
)

func Build(serviceName, remoteUrl string) error {
    tag := serviceName
    statusCode := buildImage(tag, remoteUrl)
    logger.Println(statusCode)
    if statusCode >= 300 {
        return errors.New(fmt.Sprintf("Wrong build code: %d", statusCode))
    }
    return nil
}
