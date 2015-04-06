package docker

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "strings"
)

const dockerSocket = "/var/run/docker.sock"

func getUrl(path string) string {
    return strings.Join([]string{"http://localhost", path}, "")
}

func getHttpClient() http.Client {
    tr := &http.Transport{
        Dial: func(network, addr string) (conn net.Conn, err error) {
            return net.Dial("unix", dockerSocket)
        },
    }

    return http.Client{Transport: tr}
}

func callApi(path string) http.Response {
    client := getHttpClient()
    resp, err := client.Get(strings.Join([]string{"http://localhost", path}, ""))
    if err != nil {
        log.Fatal("Unable to query API [", path, "]. Error=", err)
    }

    return *resp
}

func deleteApiCall(path string) {
    req, err := http.NewRequest("DELETE", getUrl(path), nil)
    if err != nil {
        log.Fatal("Unable to create request. ", err)
    }

    client := getHttpClient()
    _, err = client.Do(req)
    if err != nil {
        log.Fatal("Unable to delete. ", err)
    }
}

func getApiCall(path string) []interface{} {
    resp := callApi(path)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("Unable to read request body. ", err)
    }
    return parseJson(body)
}

func parseJson(data []byte) []interface{} {
    var parsedJson []interface{}
    if err := json.Unmarshal(data, &parsedJson); err != nil {
        log.Fatal("Unable to parse ", string(data))
    }
    return parsedJson
}
