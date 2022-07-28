package setters

import (
	"fmt"
	"os"

	"github.com/jeff-roche/biome/src/repos"
)

const DRAGOMAN_ENV_KEY = "from_dragoman"

// DragomanEnvironmentSetter will decrypt a secret that has been encrypted with dragoman
type DragomanEnvironmentSetter struct {
	Encrypted string                // The dragoman encrypted string
	EnvKey    string                // The environment variable to be set
	repo      repos.DragomanRepoIfc // The repo that handles decrypting with dragoman
}

// NewDragomanEnvironmentSetter is the builder function for DragomanEnvironmentSetter
func NewDragomanEnvironmentSetter(key string, val string) (*DragomanEnvironmentSetter, error) {
	dragomanRepo, err := repos.NewDragomanRepo()
	if err != nil {
		return nil, err
	}

	return &DragomanEnvironmentSetter{
		EnvKey:    key,
		Encrypted: val,
		repo:      dragomanRepo,
	}, nil
}

func (s DragomanEnvironmentSetter) SetEnv() (string, error) {
	dec, err := s.repo.Decrypt(s.Encrypted)
	if err != nil {
		return "", err
	}

	if s.EnvKey == "" {
		return "", fmt.Errorf("no environment key specified")
	}

	return dec, os.Setenv(s.EnvKey, dec)
}
