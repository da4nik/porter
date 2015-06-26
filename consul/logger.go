package consul

import (
    "log"
    "os"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[consul] ", log.Lshortfile|log.LstdFlags)
}
