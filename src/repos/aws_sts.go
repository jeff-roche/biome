package repos

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/jeff-roche/biome/src/lib/types"
)

var awsProfileName string = ""

type AwsStsRepositoryIfc interface {
	ConfigureSession(profile string) (*types.AwsEnvConfig, error)
	SetAwsEnvs(*types.AwsEnvConfig)
}

// AwsStsRepository handles setting up an AWS Session
type AwsStsRepository struct{}

func NewAwsStsRepository() *AwsStsRepository {
	return &AwsStsRepository{}
}

//
func (repo AwsStsRepository) SetAwsEnvs(cfg *types.AwsEnvConfig) {
	os.Setenv("AWS_ACCESS_KEY_ID", cfg.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", cfg.SessionToken)
	os.Setenv("AWS_DEFAULT_REGION", cfg.DefaultRegion)
}

// Configure Session will setup an AWS Session and return the needed values for the environment variables
func (repo AwsStsRepository) ConfigureSession(profile string) (*types.AwsEnvConfig, error) {
	// Load the profile
	prof, err := repo.loadProfile(profile)
	if err != nil {
		return nil, err
	}

	cfg, err := repo.loadAwsConfig(&prof)
	if err != nil {
		return nil, err
	}

	creds, err := repo.setupAwsSession(&cfg, &prof)
	if err != nil {
		return nil, err
	}

	return &types.AwsEnvConfig{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
		DefaultRegion:   prof.Region,
	}, nil
}

func (repo AwsStsRepository) loadProfile(profile string) (config.SharedConfig, error) {
	profileCfg, err := config.LoadSharedConfigProfile(context.TODO(), profile)
	if err != nil {
		return config.SharedConfig{}, err
	}

	return profileCfg, nil
}

func (repo AwsStsRepository) loadAwsConfig(profCfg *config.SharedConfig) (aws.Config, error) {
	awsProfileName = profCfg.Profile

	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profCfg.Profile),
		config.WithDefaultRegion(profCfg.Region),
		config.WithAssumeRoleCredentialOptions(
			func(aro *stscreds.AssumeRoleOptions) {
				if profCfg.MFASerial != "" {
					aro.SerialNumber = &profCfg.MFASerial
					aro.TokenProvider = CustomStdinTokenProvider
				}
			},
		),
	)
}

func (repo AwsStsRepository) setupAwsSession(cfg *aws.Config, profile *config.SharedConfig) (aws.Credentials, error) {
	// Setup an STS client
	client := sts.NewFromConfig(*cfg)

	// Get the credentials using the role
	cred_provider := stscreds.NewAssumeRoleProvider(client, profile.RoleARN)

	// Generate the temp credentials
	creds, err := cred_provider.Retrieve(context.TODO())
	if err != nil {
		return aws.Credentials{}, nil
	}

	cfg.Credentials = aws.NewCredentialsCache(cred_provider)

	return creds, nil
}

func CustomStdinTokenProvider() (string, error) {
	var v string
	fmt.Printf("MFA token for AWS profile '%s': ", awsProfileName)
	_, err := fmt.Scanln(&v)

	return v, err
}
