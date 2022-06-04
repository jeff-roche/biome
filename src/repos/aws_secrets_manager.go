package repos

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// The interface for the AWS SDK Secrets Manager Service
type awsSecretsManagerServiceIfc interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// The interface for the Secrets Manager Repository
type SecretsManagerIfc interface {
	GetSecretString(string) (string, error)
}

// The Secrets Manager Repository for proxying requests to secrets manager
type SecretsManager struct {
	client awsSecretsManagerServiceIfc
}

// NewSecretsManagerRepo builds the SecretsManagerRepository and its dependencies
func NewSecretsManagerRepo() (*SecretsManager, error) {
	// Setup the secrets manager client
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithDefaultRegion(os.Getenv("AWS_DEFAULT_REGION")),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	return &SecretsManager{
		client: secretsmanager.NewFromConfig(cfg),
	}, nil
}

// GetSecretString will pull the secret from Secrets Manager and return the string value
func (smrepo SecretsManager) GetSecretString(arn string) (string, error) {
	response, err := smrepo.client.GetSecretValue(
		context.TODO(),
		&secretsmanager.GetSecretValueInput{
			SecretId: &arn,
		},
	)

	if err != nil {
		return "", fmt.Errorf("unable to retreive '%s' from secrets manager: %v", arn, err)
	}

	return *response.SecretString, nil
}
