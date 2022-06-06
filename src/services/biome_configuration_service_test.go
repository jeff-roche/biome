package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeff-roche/biome/src/lib/types"
	"github.com/jeff-roche/biome/src/repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBiomeConfigurationService(t *testing.T) {
	sourceFilePath := "aFilePath"
	biomeName := "myBiome"
	testEnv := "MY_TEST_ENV"

	getTestBiome := func() types.Biome {
		return types.Biome{
			SourceFile: sourceFilePath,
			Config: &types.BiomeConfig{
				Name:       biomeName,
				AwsProfile: "myProfile",
				Environment: map[string]interface{}{
					testEnv: "my_test_env_var",
				},
			},
		}
	}

	t.Run("LoadBiomeFromDefaults", func(t *testing.T) {

		t.Run("should load and save the biome", func(t *testing.T) {
			// Assemble
			mockRepo := new(repos.MockBiomeFileParser)

			testBiome := getTestBiome()

			mockRepo.On("FindBiome", biomeName, mock.Anything).Return(&testBiome, nil)
			testSvc := &BiomeConfigurationService{
				configFileRepo: mockRepo,
			}

			// Act
			err := testSvc.LoadBiomeFromDefaults(biomeName)

			// Assert
			assert.Nil(t, err)
			assert.True(t, assert.ObjectsAreEqual(testBiome.Config, testSvc.ActiveBiome.Config))
			assert.Equal(t, testBiome.SourceFile, testSvc.ActiveBiome.SourceFile)
		})

		t.Run("should report an error if the biome can not be found", func(t *testing.T) {
			// Assemble
			mockRepo := new(repos.MockBiomeFileParser)
			testSvc := &BiomeConfigurationService{
				configFileRepo: mockRepo,
			}
			testErr := fmt.Errorf("biome not found")

			mockRepo.On("FindBiome", biomeName, mock.Anything).Return(&types.Biome{}, testErr)

			// Act
			err := testSvc.LoadBiomeFromDefaults(biomeName)

			// Assert
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, testErr)
			assert.Nil(t, testSvc.ActiveBiome)
		})

	})

	t.Run("LoadBiomeFromFile", func(t *testing.T) {
		t.Run("should load and save the biome", func(t *testing.T) {
			// Assemble
			mockRepo := new(repos.MockBiomeFileParser)
			testSvc := &BiomeConfigurationService{
				configFileRepo: mockRepo,
			}
			testBiome := getTestBiome()

			mockRepo.On("FindBiome", biomeName, mock.Anything).Return(&testBiome, nil)

			// Act
			err := testSvc.LoadBiomeFromFile(biomeName, sourceFilePath)

			// Assert
			assert.Nil(t, err)
			assert.True(t, assert.ObjectsAreEqual(testBiome.Config, testSvc.ActiveBiome.Config))
			assert.Equal(t, testBiome.SourceFile, testSvc.ActiveBiome.SourceFile)
		})

		t.Run("should report an error if the biome can not be found", func(t *testing.T) {
			// Assemble
			mockRepo := new(repos.MockBiomeFileParser)
			testSvc := &BiomeConfigurationService{
				configFileRepo: mockRepo,
			}
			testErr := fmt.Errorf("biome not found")

			mockRepo.On("FindBiome", biomeName, mock.Anything).Return(&types.Biome{}, testErr)

			// Act
			err := testSvc.LoadBiomeFromFile(biomeName, sourceFilePath)

			// Assert
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, testErr)
			assert.Nil(t, testSvc.ActiveBiome)
		})
	})

	t.Run("ActivateBiome", func(t *testing.T) {

		t.Run("should set the environment variable", func(t *testing.T) {
			// Assemble
			b := getTestBiome()

			testSvc := &BiomeConfigurationService{
				ActiveBiome: &b,
			}

			b.Config.AwsProfile = ""

			t.Cleanup(func() {
				os.Unsetenv(testEnv)
			})

			// Act
			err := testSvc.ActivateBiome()

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, b.Config.Environment[testEnv], os.Getenv(testEnv))
		})

		t.Run("should load the AWS environment", func(t *testing.T) {
			// Assemble
			b := getTestBiome()
			b.Config.Environment = map[string]interface{}{}
			mockRepo := repos.MockAwsStsRepository{}

			testSvc := &BiomeConfigurationService{
				ActiveBiome: &b,
				awsStsRepo:  &mockRepo,
			}

			mockRepo.On("ConfigureSession", b.Config.AwsProfile).Return(&types.AwsEnvConfig{}, nil)
			mockRepo.On("SetAwsEnvs", mock.Anything).Return()

			// Act
			err := testSvc.ActivateBiome()

			// Assert
			assert.Nil(t, err)
		})

		t.Run("should report an error if the biome is not loaded", func(t *testing.T) {
			// Assemble
			testSvc := &BiomeConfigurationService{}

			// Act
			err := testSvc.ActivateBiome()

			// Assert
			assert.NotNil(t, err)
		})

		t.Run("should report an error if the AWS profile can not be loaded", func(t *testing.T) {
			// Assemble
			b := getTestBiome()
			mockRepo := repos.MockAwsStsRepository{}
			testError := fmt.Errorf("my dummy error")
			mockRepo.On("ConfigureSession", b.Config.AwsProfile).Return(&types.AwsEnvConfig{}, testError)
			mockRepo.On("SetAwsEnvs", mock.Anything).Return()

			testSvc := &BiomeConfigurationService{
				ActiveBiome: &b,
				awsStsRepo:  &mockRepo,
			}

			// Act
			err := testSvc.ActivateBiome()

			// Assert
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, testError.Error())
		})

	})
}
