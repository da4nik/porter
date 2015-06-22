package docker

import (
    "fmt"
    "io/ioutil"
)

const (
    pathBuildImage  = "/build"          // POST
    pathPushImage   = "/images/%s/push" // POST
    pathTagImage    = "/images/%s/tag"  // POST
    pathCreateImage = "/images/create"  // POST
)

func buildImage(tag, remoteUrl string) int {

    resp, err := apiCall("POST", pathBuildImage, query{"t": {tag}, "remote": {remoteUrl}}, nil)
    if err != nil {
        logger.Fatal("Unable to create image. ", err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Fatal(err)
    }
    logger.Println(string(body))
    _ = body
    return resp.StatusCode
}

func createImage(registry, name, tag string) int {
    fromImage := fmt.Sprintf("%s/%s", registry, name)
    resp, err := apiCall("POST", pathCreateImage, query{"tag": {tag}, "fromImage": {fromImage}}, nil)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Fatal(err)
    }
    _ = body
    logger.Println(string(body))
    return resp.StatusCode
}

func tagImage(registry, name, tag string) int {
    path := fmt.Sprintf(pathTagImage, name+":"+tag)
    repo := fmt.Sprintf("%s/%s", registry, name)
    logger.Println(path)
    resp, err := apiCall("POST", path, query{"tag": {tag}, "repo": {repo}, "force": {"1"}}, nil)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Fatal(err)
    }
    _ = body
    logger.Println(string(body))
    return resp.StatusCode
}

func pushImage(registry, name, tag string) int {
    if len(registry) > 0 {
        name = fmt.Sprintf("%s/%s", registry, name)
    }
    path := fmt.Sprintf(pathPushImage, name)
    logger.Println(path)
    resp, err := apiCall("POST", path, query{"tag": {tag}}, nil)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Fatal(err)
    }
    _ = body
    logger.Println(string(body))
    return resp.StatusCode
}
