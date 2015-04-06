package consul

import (
    "strings"
)

const (
    consulUrl = "http://localhost:8500/v1/kv/"
)

func getServiceConfigKey(service string) string {
    return strings.Join([]string{"services", service, "config"}, "/")
}
