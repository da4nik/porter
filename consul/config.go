package consul

import (
  "log"
  "net/http"
  "io/ioutil"
//  "encoding/base64"
  "encoding/json"
)

type kvItem struct {
  CreateIndex int
  ModifyIndex int
  LockIndex int
  Key string
  Flags int
  Value []byte
}

type ServiceConfig struct {
  Volumes string
  Ports string
  Env string
}

func GetServiceConfig(serviceName string) *ServiceConfig {
  resp, err := http.Get(consulUrl + getServiceConfigKey(serviceName))
  if err != nil {
    log.Fatal("Unable to get ", serviceName, " config. ", err)
  }

  if resp.StatusCode != 200 {
    log.Printf("Config for service %s not found.\n", serviceName)
    return nil
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal("Unable to read config body. ", err)
  }

  var parsedJson []kvItem
  if err := json.Unmarshal(body, &parsedJson); err != nil {
    log.Fatal("Unable to parse JSON. ", string(body))
  }

  var serviceConfig ServiceConfig
  if err := json.Unmarshal(parsedJson[0].Value, &serviceConfig); err != nil {
    log.Fatal("Unable to parse service config. ", err)
  }

  return &serviceConfig
}
