package main

import (
  "log"
  "os"
  "strings"

  "github.com/da4nik/porter/docker"
)

func main() {
  if len(os.Args) < 2 {
    log.Fatal("Command not found")
  }

  switch strings.ToLower(os.Args[1]) {
  case "cleanup":
    docker.Cleanup()
  case "api":
//    docker.CallApi("/images/json")
  }

  log.Println("Success")
}
