package setters

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSecretsManagerRepo struct {
	mock.Mock
}

func (r mockSecretsManagerRepo) GetSecretString(val string) (string, error) {
	args := r.Called(val)
	return args.String(0), args.Error(1)
}

func TestSecretsManagerSetterBuilder(t *testing.T) {
	t.Run("should set all the keys specified", func(t *testing.T) {
		// Assemble
		envKey := "MY_ENV_VAR"
		configKeys := make(map[string]interface{})
		configKeys[SECRETS_MANAGER_ENV_ARN_KEY] = "myArn"
		configKeys[SECRETS_MANAGER_ENV_JSON_KEY] = "myJsonKey"

		// Act
		setter, err := NewSecretsManagerEnvironmentSetter(
			envKey,
			configKeys,
		)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, envKey, setter.EnvKey)
		assert.Equal(t, configKeys[SECRETS_MANAGER_ENV_ARN_KEY], setter.ARN)
		assert.Equal(t, configKeys[SECRETS_MANAGER_ENV_JSON_KEY], setter.SecretKey)
		assert.NotNil(t, setter.repo)
	})

	t.Run("should not set keys that are not specified", func(t *testing.T) {
		// Assemble
		envKey := "MY_ENV_VAR"
		configKeys := make(map[string]interface{})

		// Act
		setter, err := NewSecretsManagerEnvironmentSetter(
			envKey,
			configKeys,
		)

		// Assert
		assert.Nil(t, err)
		assert.Empty(t, setter.ARN)
		assert.Empty(t, setter.SecretKey)
		assert.NotNil(t, setter.repo)
	})
}

func TestSecretsManagerSetter(t *testing.T) {
	testARN := "myARN"
	testJSONKey := "mykey"
	testJSONValue := "myvalue"
	testJSON := fmt.Sprintf(`{"%s": "%s"}`, testJSONKey, testJSONValue)
	testEnv := "BIOME_TEST_ENV_VAR"

	// Helper test builder function for setup
	getTestSetter := func(mockRepo *mockSecretsManagerRepo) *SecretsManagerEnvironmentSetter {
		return &SecretsManagerEnvironmentSetter{
			EnvKey:    testEnv,
			ARN:       testARN,
			SecretKey: testJSONKey,
			repo:      mockRepo,
		}
	}

	t.Run("should set the env var from the returned JSON", func(t *testing.T) {
		// Assemble
		mockRepo := &mockSecretsManagerRepo{}
		mockRepo.On("GetSecretString", testARN).Return(testJSON, nil)
		setter := getTestSetter(mockRepo)

		t.Cleanup(func() {
			os.Unsetenv(testEnv)
		})

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, testJSONValue, val)
		assert.Equal(t, testJSONValue, os.Getenv(testEnv))
	})

	t.Run("should report an error if no ARN is specified", func(t *testing.T) {
		// Assemble
		setter := getTestSetter(nil)
		setter.ARN = ""

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, val, "")
		assert.ErrorContains(t, err, "ARN")
	})

	t.Run("should report an error getting the secret", func(t *testing.T) {
		// Assemble
		customErr := "unable to get secret"
		mockRepo := &mockSecretsManagerRepo{}
		mockRepo.On("GetSecretString", testARN).Return("", fmt.Errorf(customErr))
		setter := getTestSetter(mockRepo)

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.Equal(t, val, "")
		assert.ErrorContains(t, err, customErr)
	})

	t.Run("should report an error if no secret is returned", func(t *testing.T) {
		// Assemble
		mockRepo := &mockSecretsManagerRepo{}
		mockRepo.On("GetSecretString", testARN).Return("", nil)
		setter := getTestSetter(mockRepo)

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, val, "")
	})

	t.Run("should report an error if no JSON key is provided", func(t *testing.T) {
		// Assemble
		mockRepo := &mockSecretsManagerRepo{}
		mockRepo.On("GetSecretString", testARN).Return(testJSON, nil)
		setter := getTestSetter(mockRepo)
		setter.SecretKey = ""

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, val, "")
		assert.ErrorContains(t, err, "JSON")
	})

	t.Run("should report an error if invalid JSON is returned", func(t *testing.T) {
		// Assemble
		mockRepo := &mockSecretsManagerRepo{}
		mockRepo.On("GetSecretString", testARN).Return("I'm Not Valid }", nil)
		setter := getTestSetter(mockRepo)

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, val, "")
		assert.ErrorContains(t, err, "unable to parse secret")
	})

	t.Run("should report an error if the specified JSON key does not exist", func(t *testing.T) {
		// Assemble
		mockRepo := &mockSecretsManagerRepo{}
		mockRepo.On("GetSecretString", testARN).Return(testJSON, nil)
		setter := getTestSetter(mockRepo)
		setter.SecretKey = "invalidKey"

		// Act
		val, err := setter.SetEnv()

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, val, "")
		assert.ErrorContains(t, err, "does not contain JSON key")
	})
}
