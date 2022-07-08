package setters

import (
	"fmt"
	"os"
)

func GetEnvironmentSetter(key string, node interface{}) (EnvironmentSetter, error) {
	var setter EnvironmentSetter
	switch node.(type) {
	case map[string]interface{}: // Complex keys
		var err error
		setter, err = getComplexSetter(key, node.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
	default: // basic types
		setter = NewBasicEnvironmentSetter(key, node)
	}

	return setter, nil
}

func getComplexSetter(key string, node map[string]interface{}) (EnvironmentSetter, error) {
	// Secrets Manager Secret
	if _, exists := node[SECRETS_MANAGER_ENV_ARN_KEY]; exists {

		return NewSecretsManagerEnvironmentSetter(key, node)
	}

	// Dragoman Encrypted Secret
	if val, exists := node[DRAGOMAN_ENV_KEY]; exists {
		return NewDragomanEnvironmentSetter(key, val.(string))
	}

	// CLI Input
	if val, exists := node[CLI_ENVIRONMENT_SETTER_KEY]; exists && val.(bool) {
		var isSecret bool

		if _, exists := node[CLI_ENVIRONMENT_SECRET_SETTER_KEY]; exists {
			isSecret = node[CLI_ENVIRONMENT_SECRET_SETTER_KEY].(bool)
		} else {
			isSecret = false
		}

		return NewCLIEnvironmentSetter(key, os.Stdin, isSecret)
	}

	return nil, fmt.Errorf("unkown environment config for variable '%s'", key)
}
