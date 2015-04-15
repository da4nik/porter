package consul

import (
    "strings"
    "net/http"
    "log"
    "io/ioutil"
)

const (
    urlConsul = "http://localhost:8500/v1"
)

func getServiceConfigKey(service string) string {
    return strings.Join([]string{"services", service, "config"}, "/")
}

func apiCall(path string) []byte {
    resp, err := http.Get(urlConsul + path)
    if err != nil {
        log.Fatal("Unable to get call consul api. ", err)
    }

    if resp.StatusCode != 200 {
        log.Fatal("Consul API returned error.", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("Unable to read config body. ", err)
    }
    return body
}
