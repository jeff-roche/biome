package setters

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jeff-roche/biome/src/repos"
)

const SECRETS_MANAGER_ENV_ARN_KEY = "secret_arn"
const SECRETS_MANAGER_ENV_JSON_KEY = "secret_json_key"

// SecretsManagerEnvironmentSetter will set an environment variable
// from a JSON secret stored in AWS Secrets Manager
type SecretsManagerEnvironmentSetter struct {
	EnvKey    string                  // The environment variable key being set
	ARN       string                  // The secrets manager secret ARN to reference
	SecretKey string                  // The JSON key in the secret to use
	repo      repos.SecretsManagerIfc // The Secrets manager client
}

// NewSecretsManagerEnvironmentSetter will generete a SM Setter
func NewSecretsManagerEnvironmentSetter(key string, subkeys map[string]interface{}) (*SecretsManagerEnvironmentSetter, error) {
	setter := &SecretsManagerEnvironmentSetter{
		EnvKey: key,
	}

	// ARN
	if val, exists := subkeys[SECRETS_MANAGER_ENV_ARN_KEY]; exists {
		setter.ARN = val.(string)
	}

	// JSON Key
	if val, exists := subkeys[SECRETS_MANAGER_ENV_JSON_KEY]; exists {
		setter.SecretKey = val.(string)
	}

	// Secrets Manager Repo
	var err error
	setter.repo, err = repos.NewSecretsManagerRepo()

	if err != nil {
		return nil, fmt.Errorf("unable to initialize the secrets manager repository: %v", err)
	}

	return setter, nil
}

func (s SecretsManagerEnvironmentSetter) SetEnv() (string, error) {
	// Get the value from Secrets Manager
	if s.ARN == "" {
		return "", NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("no Secrets Manager ARN specified. Please use '%s' to specify one", SECRETS_MANAGER_ENV_ARN_KEY),
		)
	}

	secret, err := s.repo.GetSecretString(s.ARN)
	if err != nil {
		return "", NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("unable to load secret '%s': %v", s.ARN, err),
		)
	}

	// Validate the output
	if secret == "" {
		return "", NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("only string based secrets are currently supported"),
		)
	}

	// Get value from the JSON key
	if s.SecretKey == "" {
		return "", NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("no JSON key for the secret was specified. Please use '%s' to specify one", SECRETS_MANAGER_ENV_JSON_KEY),
		)
	}

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(secret), &jsonData); err != nil {
		return "", NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("unable to parse secret '%s' JSON: %v", s.ARN, err),
		)
	}

	// Check and set the key
	val, exists := jsonData[s.SecretKey]

	if !exists {
		return "", NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("secret '%s' does not contain JSON key '%s'", s.ARN, s.SecretKey),
		)
	}

	return val, os.Setenv(s.EnvKey, val)
}

type SecretsManagerEnvironmentSetterError struct {
	varName string
	value   string
}

func (e SecretsManagerEnvironmentSetterError) Error() string {
	return fmt.Sprintf("error setting env var '%s' from Secrets Manager value: %s", e.varName, e.value)
}

func NewSecretsManagerEnvironmentSetterError(variable string, err string) error {
	return SecretsManagerEnvironmentSetterError{
		varName: variable,
		value:   err,
	}
}
