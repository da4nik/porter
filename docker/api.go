package docker

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "strings"
    "bytes"
)

const dockerSocket = "/var/run/docker.sock"

func getHttpClient() http.Client {
    tr := &http.Transport{
        Dial: func(network, addr string) (conn net.Conn, err error) {
            return net.Dial("unix", dockerSocket)
        },
    }

    return http.Client{Transport: tr}
}

func jsonGetApiCall(path string) []interface{} {
    resp, err := apiCall("GET", path, nil, nil)
    if err != nil {
        log.Fatal("Unable to complete GET request. ", err)
    }

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

func getApiUrl(path string, params []string) string {
    var url_params string
    if len(params) > 0 {
        url_params = "?" + strings.Join(params, "&")
    }
    return strings.Join([]string{"http://localhost", path, url_params}, "")
}

func apiCall(method, path string, params []string, payload []byte) (resp *http.Response, error error) {
    req, err := http.NewRequest(strings.ToUpper(method), getApiUrl(path, params), bytes.NewReader(payload))
    if err != nil {
        log.Fatal("Unable to create request. ", err)
    }
    req.Header.Add("Content-type", "application/json")

    client := getHttpClient()
    return client.Do(req)
}
