package docker

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "log"
    "os/exec"
    "strings"
    "testing"
)

func shRunContainer(name string) string {
    id, err := exec.Command("docker", "run", "-d", fmt.Sprintf("--name=%s", name), "busybox", "tail", "-f", "/etc/issue").Output()
    if err != nil {
        log.Fatal("run container error: ", err)
    }
    return strings.TrimSuffix(string(id), "\n")
}

func shRemoveContainer(name string) {
    exec.Command("docker", "rm", "-f", name).Output()
}

func shInspectContainer(name, format string) string {
    data, err := exec.Command("docker", "inspect", fmt.Sprintf("--format={{%s}}", format), name).Output()
    if err != nil {
        log.Fatal("inspect container error: ", err)
    }
    return strings.TrimSuffix(string(data), "\n")
}

func containerIsRunning(name string) bool {
    inpectData := shInspectContainer(name, ".State.Running")
    return inpectData == "true"
}

func containerId(name string) string {
    return shInspectContainer(name, ".Id")
}

func containerExists(name string) bool {
    data, err := exec.Command("docker", "ps", "-a").Output()
    if err != nil {
        log.Fatal("list container error: ", err)
    }
    return strings.Contains(string(data), name)
}

func TestStopContainer(t *testing.T) {
    assert := assert.New(t)
    name := "test_container"
    shRemoveContainer(name)
    shRunContainer(name)
    code := stopContainer(name)
    assert.Equal(204, code, "wrong code")
    assert.False(containerIsRunning(name))
}

func TestRenameContainer(t *testing.T) {
    assert := assert.New(t)
    name := "test_container"
    new_name := "test_conatiner_2"
    shRemoveContainer(name)
    shRemoveContainer(new_name)
    id := shRunContainer(name)
    code := renameContainer(name, new_name)
    assert.Equal(204, code, "wrong code")
    assert.False(containerExists(name))
    assert.True(containerExists(new_name))
    _ = id
    assert.Equal(id, containerId(new_name))
    // assert.True(containerIsRunning(new_name))
}

func TestRemoveContainer(t *testing.T) {
    assert := assert.New(t)
    name := "test_container"
    shRemoveContainer(name)
    shRunContainer(name)
    stopContainer(name)
    code := removeContainer(name)
    assert.Equal(204, code, "wrong code")
    assert.False(containerExists(name))
}
