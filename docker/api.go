package docker

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net"
    "net/http"
    "net/url"
    "strings"
)

type query map[string][]string

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
        logger.Fatal("Unable to complete GET request. ", err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Fatal("Unable to read request body. ", err)
    }
    return parseJson(body)
}

func parseJson(data []byte) []interface{} {
    var parsedJson []interface{}
    if err := json.Unmarshal(data, &parsedJson); err != nil {
        logger.Fatal("Unable to parse ", string(data))
    }
    return parsedJson
}

func getApiUrl(path string, params map[string][]string) (apiUrl string) {
    apiUrl = "http://localhost" + path
    u, err := url.Parse(apiUrl)
    if err != nil {
        logger.Fatal(err)
    }
    q := u.Query()
    for k, vs := range params {
        for _, v := range vs {
            q.Set(k, v)
        }
    }
    u.RawQuery = q.Encode()
    apiUrl = u.String()
    return
}

func apiCall(method, path string, params query, payload []byte) (resp *http.Response, error error) {
    req, err := http.NewRequest(strings.ToUpper(method), getApiUrl(path, params), bytes.NewReader(payload))
    if err != nil {
        logger.Fatal("Unable to create request. ", err)
    }
    req.Header.Add("Content-type", "application/json")
    req.Header.Add("X-Registry-Auth", "0")
    client := getHttpClient()
    return client.Do(req)
}
