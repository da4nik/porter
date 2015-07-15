package consul

import (
    "fmt"
    consulapi "github.com/hashicorp/consul/api"
)

type NotUpdatedError struct {
    key string
}

func (e NotUpdatedError) Error() string {
    return fmt.Sprintf("Record '%s' already chaged", e.key)
}

type NoConfigError struct {
    key string
}

func (e NoConfigError) Error() string {
    return fmt.Sprintf("No config for key %s\n", e.key)
}

type Config interface {
    serialize() ([]byte, error)
    deserialize([]byte) error
    Key() string
    SetModifyIndex(uint64)
    GetModifyIndex() uint64
}

func LoadConfig(c Config) error {
    pair, err := api.GetKVPair(c.Key())
    if err != nil {
        return err
    }
    return fillFromKV(c, pair)
}

func fillFromKV(c Config, pair *consulapi.KVPair) error {
    if pair == nil {
        return NoConfigError{c.Key()}
    }
    c.SetModifyIndex(pair.ModifyIndex)
    return c.deserialize(pair.Value)
}

func SaveConfig(c Config) error {
    pair := new(consulapi.KVPair)
    pair.Key = c.Key()
    pair.ModifyIndex = c.GetModifyIndex()
    value, err := c.serialize()
    if err != nil {
        return err
    }
    pair.Value = value
    return api.PutKVPair(pair)
}

func DeleteConfig(c Config) error {
    return api.DeleteKvPair(c.Key())
}
