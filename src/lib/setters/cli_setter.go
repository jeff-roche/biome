package setters

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const CLI_ENVIRONMENT_SETTER_KEY = "from_cli"
const CLI_ENVIRONMENT_SECRET_SETTER_KEY = "is_secret"

// CLIEnvironmentSetter will pull an env var from the provided io.Reader (stdin)
type CLIEnvironmentSetter struct {
	Key   string
	Value string
}

func NewCLIEnvironmentSetter(key string, rd io.Reader, isSecret bool) (*CLIEnvironmentSetter, error) {
	// Get the value from the io.Reader
	fmt.Printf("%s: ", key)
	var val string
	var err error

	if isSecret {
		val, err = getSecretCliInput(rd)
	} else {
		val, err = getCliInput(rd)
	}

	if err != nil {
		return nil, err
	}

	// Save and return the setter
	return &CLIEnvironmentSetter{
		Key:   key,
		Value: val,
	}, nil
}

func (s CLIEnvironmentSetter) SetEnv() (string, error) {
	return s.Value, os.Setenv(s.Key, s.Value)
}

func getCliInput(rd io.Reader) (string, error) {
	reader := bufio.NewReader(rd)
	val, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	val = strings.Replace(val, "\n", "", -1) // Remove any newline characters

	return val, nil
}

func getSecretCliInput(rd io.Reader) (string, error) {
	byteSecret, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}

	return strings.Replace(string(byteSecret), "\n", "", -1), nil // Convert to string and remove any newline characters
}
