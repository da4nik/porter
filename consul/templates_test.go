package consul

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "os"
    "testing"
)

const (
    destName = "/tmp/1"
    newName  = "/tmp/2"
)

func newTemplate() *ServiceTemplate {
    tmpl := new(ServiceTemplate)
    tmpl.service = new(ServiceConfig)
    tmpl.service.Name = "service_1"
    tmpl.Command = fmt.Sprintf("cp %s %s", destName, newName)
    tmpl.Name = "conf_1"
    tmpl.Template = "{{with node}}nodename - {{.Node.Node}}{{end}}"
    tmpl.Destination = destName
    return tmpl
}

func TestRun(t *testing.T) {
    os.Remove(destName)
    os.Remove(newName)
    assert := assert.New(t)
    tmpl := newTemplate()
    tmpl.Run()
    _, err := os.Stat(destName)
    assert.NoError(err)
    _, err = os.Stat(newName)
    assert.NoError(err)
    rendered, err := ioutil.ReadFile(destName)
    assert.NoError(err)
    assert.Contains(string(rendered), "nodename")
}

func TestSave(t *testing.T) {
    assert := assert.New(t)
    tmpl := newTemplate()
    t.Log(tmpl.Key())
    api.DeleteKvPair(tmpl.Key())
    err := tmpl.Update()
    assert.NoError(err)
    _, err = api.GetKVPair(tmpl.Key())
    assert.NoError(err)
}
