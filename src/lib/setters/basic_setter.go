package setters

import (
	"fmt"
	"os"
)

type BasicEnvironmentSetter struct {
	Key   string
	Value string
}

func NewBasicEnvironmentSetter(key string, value interface{}) *BasicEnvironmentSetter {
	return &BasicEnvironmentSetter{
		Key:   key,
		Value: fmt.Sprint(value),
	}
}

func (s BasicEnvironmentSetter) SetEnv() (string, error) {
	return s.Value, os.Setenv(s.Key, s.Value)
}
