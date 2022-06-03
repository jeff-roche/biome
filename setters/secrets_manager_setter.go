package setters

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const SECRETS_MANAGER_ENV_ARN_KEY = "secret_arn"
const SECRETS_MANAGER_ENV_JSON_KEY = "secret_json_key"

// SecretsManagerEnvironmentSetter will set an environment variable
// from a JSON secret stored in AWS Secrets Manager
type SecretsManagerEnvironmentSetter struct {
	EnvKey    string // The environment variable key being set
	ARN       string // The secrets manager secret ARN to reference
	SecretKey string // The JSON key in the secret to use
}

// NewSecretsManagerEnvironmentSetter will generete a SM Setter
func NewSecretsManagerEnvironmentSetter(key string, config map[string]interface{}) *SecretsManagerEnvironmentSetter {
	setter := &SecretsManagerEnvironmentSetter{
		EnvKey: key,
	}

	// ARN
	if val, exists := config[SECRETS_MANAGER_ENV_ARN_KEY]; exists {
		setter.ARN = val.(string)
	}

	// JSON Key
	if val, exists := config[SECRETS_MANAGER_ENV_JSON_KEY]; exists {
		setter.SecretKey = val.(string)
	}

	return setter
}

func (s SecretsManagerEnvironmentSetter) SetEnv() error {

	// Setup the secrets manager client
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithDefaultRegion(os.Getenv("AWS_DEFAULT_REGION")),
	)
	if err != nil {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("unable to load AWS configuration: %v", err),
		)
	}

	client := sm.NewFromConfig(cfg)

	// Get the value from Secrets Manager
	if s.ARN == "" {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("no Secrets Manager ARN specified. Please use '%s' to specify one", SECRETS_MANAGER_ENV_ARN_KEY),
		)
	}

	smOut, err := client.GetSecretValue(
		context.TODO(),
		&sm.GetSecretValueInput{
			SecretId: &s.ARN,
		},
	)
	if err != nil {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("unable to load secret '%s': %v", s.ARN, err),
		)
	}

	// Validate the output
	if *smOut.SecretString == "" {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("only string based secrets are currently supported"),
		)
	}

	// Get value from the JSON key
	if s.ARN == "" {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("no JSON key for the secret was specified. Please use '%s' to specify one", SECRETS_MANAGER_ENV_JSON_KEY),
		)
	}

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(*smOut.SecretString), &jsonData); err != nil {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("unable to parse secret '%s' JSON: %v", s.ARN, err),
		)
	}

	// Check and set the key
	if val, exists := jsonData[s.SecretKey]; exists {
		os.Setenv(s.EnvKey, val)
	} else {
		return NewSecretsManagerEnvironmentSetterError(
			s.EnvKey,
			fmt.Sprintf("secret '%s' does not contain JSON key '%s'", s.ARN, s.SecretKey),
		)
	}

	return nil
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
