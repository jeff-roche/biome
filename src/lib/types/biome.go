package types

type Biome struct {
	Config     *BiomeConfig
	SourceFile string
}

type BiomeConfig struct {
	Name        string
	AwsProfile  string                 `yaml:"aws_profile"`
	Commands    []string               `yaml:"commands"`
	Environment map[string]interface{} `yaml:"environment"`
}
