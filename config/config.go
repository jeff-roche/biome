package config

import (
	"fmt"
	"log"
)

type BiomeConfig struct {
	Name string
	AwsProfile string `yaml:"aws_profile"`
	Environment map[string]interface{} `yaml:"environment"`
}

func (bc BiomeConfig) GetEnvs() map[string]string {
	newEnvs := make(map[string]string)

	for key, node := range bc.Environment {
		var val string
		switch node.(type) {
		case map[string]interface{}:
			var err error
			val, err = parseBiomeVariableConfig(node.(map[string]interface{}))
			if err != nil {
				log.Fatalf("Error processing variable '%s': %v", key, err)
			}
		default:
			val = fmt.Sprint(node)
		}

		newEnvs[key] = val
	}

	return newEnvs
}

type BiomeVariableConfig struct {
	SecretARN string
	SecretJSONKey string
}

func parseBiomeVariableConfig(cfg map[string]interface{}) (string, error) {
	bvc := BiomeVariableConfig{}
	
	// Parse out any keys, and see what we've got
	for key, val := range cfg {
		switch val.(type) {
		case string:
			parsedVal := val.(string)
			switch key {
			case "secret_arn":
				bvc.SecretARN = parsedVal
			case "secret_json_key":
				bvc.SecretJSONKey = parsedVal
			default:
				return "", fmt.Errorf("unable to process key %s: unknown key", key)
			}
		default:
			return "", fmt.Errorf("unable to process key %s with value %v", key, val)
		}
	}

	// Process a secrets manager secret
	if bvc.SecretARN != "" {
		/*smval, err := loadSecretsManagerSecret(bvc.SecretARN)
		if err != nil {
			return "", fmt.Errorf("unable to load secret '%s' from secrets manager: %v", bvc.SecretARN, err)
		}*/

		/*if bvc.SecretJSONKey != "" {
			jsonData := make(map[string]string)
			json.Unmarshal([]byte(smval), &jsonData)

			return jsonData[smval], nil
		} else {
			return smval, nil
		}*/

		return bvc.SecretARN, nil
	}

	// Could not process this variable, error time
	return "", fmt.Errorf("unknown error occured while parsing the variable")

}