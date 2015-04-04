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

func GetServiceConfig(serviceName string) kvItem {
  resp, err := http.Get(consulUrl + getServiceConfigKey(serviceName))
  if err != nil {
    log.Fatal("Unable to get ", serviceName, " config. ", err)
  }

  if resp.StatusCode != 200 {
    log.Printf("Config for service %s not found.\n", serviceName)
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal("Unable to read config body. ", err)
  }

  var parsedJson []kvItem
  if err := json.Unmarshal(body, &parsedJson); err != nil {
    log.Fatal("Unable to parse JSON. ", string(body))
  }

  log.Println("PARSED", string(parsedJson[0].Value))

  return parsedJson[0]
}
