package main

import (
    "fmt"
    "os"
)

const VERSION = "0.1.2"

func Version() {
    fmt.Fprintf(os.Stdout, "%s\n", VERSION)
}
