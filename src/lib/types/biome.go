package types

import (
	"fmt"
)

type Biome struct {
	Config     *BiomeConfig
	SourceFile string
}

type BiomeConfig struct {
	Name            string
	AwsProfile      string                 `yaml:"aws_profile"`
	Commands        []string               `yaml:"commands"`
	ExternalEnvFile string                 `yaml:"load_env"`
	Environment     map[string]interface{} `yaml:"environment"`
	Inheritance     string                 `yaml:"inherit_from"`
}

func (bc *BiomeConfig) Inherit(genepool map[string]*BiomeConfig) error {

	parents := make([]string, 0, len(genepool))
	unique := make(map[string]bool, len(genepool))
	inherits_from := bc.Inheritance

	for inherits_from != "" {
		// Does it exist?
		if _, exists := genepool[inherits_from]; !exists {
			return fmt.Errorf("can not find inherited biome %s", inherits_from)
		}

		// Ciruclar inheritance check
		if _, exists := unique[inherits_from]; exists {
			return fmt.Errorf("circular biome inheritance found, %s inherited cyclicly", inherits_from)
		}

		unique[inherits_from] = true

		parents = append(parents, inherits_from)
		inherits_from = genepool[inherits_from].Inheritance
	}

	// Now do the inheritance
	for _, i := range parents {
		biome := genepool[i]
		// AWS Profile (if one hasn't been set)
		if bc.AwsProfile == "" {
			bc.AwsProfile = biome.AwsProfile
		}

		// Load in any external env files (if one isn't specified yet)
		if bc.ExternalEnvFile == "" {
			bc.ExternalEnvFile = biome.ExternalEnvFile
		}

		// Envs (only if they don't already exist)
		for env, val := range biome.Environment {
			if _, exists := bc.Environment[env]; !exists {
				bc.Environment[env] = val
			}
		}

		// Commands (prepend to the commands array)
		bc.Commands = append(biome.Commands, bc.Commands...)
	}

	return nil
}
