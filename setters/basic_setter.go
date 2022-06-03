package setters

import "os"

type BasicEnvironmentSetter struct {
	Key   string
	Value string
}

func NewBasicEnvironmentSetter(key string, value string) *BasicEnvironmentSetter {
	return &BasicEnvironmentSetter{
		Key:   key,
		Value: value,
	}
}

func (s BasicEnvironmentSetter) SetEnv() error {
	return os.Setenv(s.Key, s.Value)
}
