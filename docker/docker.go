package docker

import (
  "os/exec"
  "log"
  "strings"
)

const dockerCommand = "docker"

func init() {
  _, err := exec.LookPath("docker")
  if err != nil {
    log.Fatal("Unable to find docker executable")
  }
}

func exitedContainers() string {
  out, err := exec.Command(dockerCommand, "ps", "-f", "status=exited", "-q").Output()
  if err != nil {
    log.Fatal("Unable to get exited containers ids. ", err)
  }
  return strings.Replace(string(out), "\n", " ", -1)
}

func untaggedImages() string {
  out, err := exec.Command(dockerCommand, "images", "-q", "-f", "dangling=true").Output()
  if err != nil {
    log.Fatal("Unable to get untagged containers ids. ", err)
  }
  return strings.Replace(string(out), "\n", " ", -1)
}

func execBashedCommand(command string) {
  if err := exec.Command("/bin/sh", "-c", command).Run(); err != nil {
    log.Fatal("Unable to execute command [", command, "] ", err);
  }
}

// Remove stopped containers and <none> images
func Cleanup() {
  containers := exitedContainers()
  if len(containers) > 0 {
    command := strings.Join([]string{"docker rm", containers}, " ")
    execBashedCommand(command)
    log.Println("Exited containers cleared.")
  }

  images :=untaggedImages()
  if len(images) > 0 {
    command := strings.Join([]string{"docker rmi", images}, " ")
    execBashedCommand(command)
    log.Println("Untagged containers cleared.")
  }

  log.Println("Cleanup complete.")
}
