package utils

import (
    "encoding/json"
    "log"
)

func ParseJson(data []byte) []interface{} {
    var parsedJson []interface{}
    if err := json.Unmarshal(data, &parsedJson); err != nil {
        log.Fatal("Unable to parse JSON. ", string(data))
    }
    return parsedJson
}
