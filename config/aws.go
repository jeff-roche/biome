package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// ConfigureAwsEnvironment will load the profile and setup the aws environment variables needed for AWS command execution
func ConfigureAwsEnvironment(profile string) error {
	prof, err := loadAwsProfile(profile)
	if err != nil {
		return fmt.Errorf("unable to load AWS Profile '%s': %v", profile, err)
	}

	cfg, err := loadAwsConfig(&prof)
	if err != nil {
		return fmt.Errorf("unable to load AWS Config in '%s': %v", prof.Region, err)
	}

	creds, err := setupAwsSession(&cfg, &prof)
	if err != nil {
		return fmt.Errorf("unable to setup the AWS Session for profile '%s': %v", prof.Profile, err)
	}

	setupAwsEnv(&creds, prof.Region)

	return nil
}

// loadAwsProfile will load data from the profile in ~/.aws/credentials
func loadAwsProfile(profile string) (config.SharedConfig, error) {
	profilecfg, err := config.LoadSharedConfigProfile(context.TODO(), profile)
	if err != nil {
		return config.SharedConfig{}, err
	}

	return profilecfg, nil
}

// loadAwsConfig will setup the configuration object to communicate with STS
func loadAwsConfig(profile *config.SharedConfig) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile.Profile),
		config.WithDefaultRegion(profile.Region),
	)
	
	if err != nil {
		return aws.Config{}, err
	}
	
	return cfg, nil
}

// setupAwsSession will generate the temp credentials from the config and profile
func setupAwsSession(cfg *aws.Config, profile *config.SharedConfig) (aws.Credentials, error) {
	// Setup an STS client
	client := sts.NewFromConfig(*cfg)

	// Get the credentials using the role
	cred_provider := stscreds.NewAssumeRoleProvider(client, profile.RoleARN)

	// Generate the temp credentials
	creds, err := cred_provider.Retrieve(context.TODO())
	if err != nil {
		return aws.Credentials{}, nil
	}

	return creds, nil
}

// setupAwsEnv will use the credentials and region to set the standard AWS Environment Variables
func setupAwsEnv(creds *aws.Credentials, region string) {
	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
	os.Setenv("AWS_DEFAULT_REGION", region)
}

// loadSecretsManagerSecret will pull in a secret from AWS Secrets Manager
func loadSecretsManagerSecret(id string) (string, error) {
	var region = os.Getenv("AWS_DEFAULT_REGION")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithDefaultRegion(region),
	)
	if err != nil {
		return "", err
	}

	client := sm.NewFromConfig(cfg)

	out, err := client.GetSecretValue(
		context.TODO(),
		&sm.GetSecretValueInput{
			SecretId: &id,
		},
	)
	if err != nil {
		return "", err
	}

	if *out.SecretString == "" {
		return "", fmt.Errorf("only string secrets are supported")
	}

	return *out.SecretString, nil
}