package docker

import (
	"flag"
	"log"
	"net/http"

	"bytes"
	"encoding/json"
	"github.com/da4nik/porter/consul"
	"io/ioutil"
)

const (
	pathCreateContainer = "/containers/create" // POST
)

func createContainer(serviceName string) {

}

func runContainer(serviceName string) {

}

func Run() {
	if flag.Arg(1) == "" {
		log.Fatal("No service name given")
	}

	serviceName := flag.Arg(1)
	serviceConfig := consul.GetServiceConfig(serviceName)

	parameters := map[string]interface{}{
		"Image": serviceName,
	}

	if serviceConfig != nil {
		// set params
		log.Println("SETTING CONFIGS")
		log.Println(serviceConfig.Volumes)
	}

	params_json, _ := json.Marshal(parameters)
	log.Println("Container params ", string(params_json))

	req, err := http.NewRequest("POST", getUrl(pathCreateContainer)+"?name="+serviceName, bytes.NewReader(params_json))
	if err != nil {
		log.Fatal("Unable to create request. ", err)
	}
	req.Header.Add("Content-type", "application/json")

	client := getHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Unable to create container. ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read request body. ", err)
	}
	log.Println(resp)
	log.Println(string(body))
}
