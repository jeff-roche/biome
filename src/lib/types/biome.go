package types

type Biome struct {
	Config     *BiomeConfig
	SourceFile string
}

type BiomeConfig struct {
	Name        string
	AwsProfile  string                 `yaml:"aws_profile"`
	Environment map[string]interface{} `yaml:"environment"`
}
