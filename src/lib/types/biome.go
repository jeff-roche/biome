package types

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
}
