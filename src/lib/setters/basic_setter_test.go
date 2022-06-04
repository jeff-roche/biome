package setters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicSetter(t *testing.T) {
	testEnvKey := "FOOBAR_TEST_KEY"
	t.Run("should set the key and value", func(t *testing.T) {
		testVal := "BAZ"
		s := NewBasicEnvironmentSetter(testEnvKey, testVal)

		assert.Equal(t, testEnvKey, s.Key)
		assert.Equal(t, testVal, s.Value)
	})

	t.Run("should convert the value to a string", func(t *testing.T) {
		testVal := 12345
		s := NewBasicEnvironmentSetter(testEnvKey, testVal)

		assert.Equal(t, testEnvKey, s.Key)
		assert.Equal(t, "12345", s.Value)
	})
}

func TestBasicSetter(t *testing.T) {
	testEnvKey := "FOOBAR_TEST_KEY"
	t.Run("should set the env", func(t *testing.T) {
		testVal := "BAZ"
		s := NewBasicEnvironmentSetter(testEnvKey, testVal)
		err := s.SetEnv()

		assert.Nil(t, err)
		assert.Equal(t, os.Getenv(testEnvKey), testVal)
	})

	os.Unsetenv(testEnvKey)
}
