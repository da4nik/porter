package docker

import (
    "log"
    "os"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[docker] ", log.Lshortfile|log.LstdFlags)
}
