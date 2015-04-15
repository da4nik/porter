package consul

import (
    "flag"
    "log"
    "encoding/json"
    "os"
    "fmt"
)

const (
    pathServicesHosts = "/catalog/service/"
)

type ServiceNode struct {
    Node string
    Address string
    ServiceID string
    ServiceName string
    ServiceTags []string
    ServiceAddress string
    ServicePort int
}

func Addresses() {
    if flag.Arg(1) == "" {
        log.Fatal("No service name given")
    }

    serviceName := flag.Arg(1)

    body := apiCall(pathServicesHosts + serviceName)

    var nodes []ServiceNode
    if err := json.Unmarshal(body, &nodes); err != nil {
        log.Fatal("Unable to parse JSON. ", string(body))
    }

    for _, node := range nodes {
        host := node.ServiceAddress
        if len(host) == 0 { host = node.Address }
        if len(host) == 0 { continue }
        if node.ServicePort > 0 { host = fmt.Sprintf("%s:%d", host, node.ServicePort) }
        fmt.Fprint(os.Stdout, host, "\n")
    }
}

