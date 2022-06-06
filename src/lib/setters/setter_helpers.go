package setters

import (
	"fmt"
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

	return nil, fmt.Errorf("unkown environment config for variable '%s'", key)
}
