package consul

import (
    "encoding/json"
    "fmt"
    "github.com/hashicorp/consul/command"
    "github.com/mitchellh/cli"
    "io/ioutil"
    "os"
    "os/signal"
    "regexp"
    "strings"
    "syscall"
)

func preparePath(path string) string {
    re := regexp.MustCompile("\\s")
    return re.ReplaceAllString(path, "_")
}

type ServiceTemplate struct {
    Name        string
    Template    string
    Destination string
    Command     string
    ModifyIndex uint64
    service     *ServiceConfig
}

func (s *ServiceTemplate) serialize() (result []byte, err error) {
    return json.Marshal(s)
}

func (s *ServiceTemplate) deserialize(data []byte) error {
    return json.Unmarshal(data, s)
}

func (t *ServiceTemplate) Key() string {
    return fmt.Sprintf("%s/templates/%s", t.service.Key(), t.Name)
}

func (s *ServiceTemplate) GetModifyIndex() uint64 {
    return s.ModifyIndex
}

func (s *ServiceTemplate) SetModifyIndex(index uint64) {
    s.ModifyIndex = index
}

func (s *ServiceTemplate) Update() error {
    return SaveConfig(s)
}

func (c *ServiceTemplate) templateFile() (f *os.File, err error) {
    prefix := fmt.Sprintf("config-%s-%s", c.service.Name, c.Name)
    prefix = preparePath(prefix)
    f, err = ioutil.TempFile("/tmp", prefix)
    if err != nil {
        return
    }
    defer f.Close()
    _, err = f.WriteString(c.Template)
    if err != nil {
        return
    }
    return
}

func removeTemplateFile(f *os.File) {
    os.Remove(f.Name())
}

func (c *ServiceTemplate) FullCommand(templateFile *os.File) []string {
    templateParams := []string{
        templateFile.Name(),
        c.Destination,
    }
    if c.Command != "" {
        templateParams = append(templateParams, c.Command)
    }
    result := []string{
        "consul-template",
        "-consul 127.0.0.1:8500",
        fmt.Sprintf("-template \"%s\"", strings.Join(templateParams, ":")),
        "-once",
    }
    return result
}

func makeShutdownCh() <-chan struct{} {
    resultCh := make(chan struct{})

    signalCh := make(chan os.Signal, 4)
    signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
    go func() {
        for {
            <-signalCh
            resultCh <- struct{}{}
        }
    }()

    return resultCh
}

func (c *ServiceTemplate) Run() {
    f, err := c.templateFile()
    defer removeTemplateFile(f)
    if err != nil {
        logger.Fatal(err)
    }
    ui := &cli.BasicUi{Writer: os.Stdout}
    ec := &command.ExecCommand{
        ShutdownCh: makeShutdownCh(),
        Ui:         ui,
    }
    ec.Run(c.FullCommand(f))
}
