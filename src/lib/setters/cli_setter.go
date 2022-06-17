package setters

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const CLI_ENVIRONMENT_SETTER_KEY = "from_cli"

// CLIEnvironmentSetter will pull an env var from the provided io.Reader (stdin)
type CLIEnvironmentSetter struct {
	Key   string
	Value string
}

func NewCLIEnvironmentSetter(key string, rd io.Reader) (*CLIEnvironmentSetter, error) {
	// Get the value from the io.Reader
	fmt.Printf("%s: ", key)
	reader := bufio.NewReader(rd)
	val, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	val = strings.Replace(val, "\n", "", -1) // Remove any newline characters

	// Save and return the setter
	return &CLIEnvironmentSetter{
		Key:   key,
		Value: val,
	}, nil
}

func (s CLIEnvironmentSetter) SetEnv() error {
	return os.Setenv(s.Key, s.Value)
}
