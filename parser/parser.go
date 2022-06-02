package parser

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/jeff-roche/biome/config"

	"gopkg.in/yaml.v3"
)

var biomeConfigDefaultFileNames = []string{".biome.yaml", ".biome.yml"}

// BiomeParser is the parent struct that handles loading and parsing biome configs
type BiomeParser struct {
	LoadedBiome          *config.BiomeConfig
	BiomeConfigFileNames []string
}

// NewBiomeParser will generate a BiomeParser with the default configuration
func NewBiomeParser() *BiomeParser {
	return &BiomeParser{
		LoadedBiome:          nil,
		BiomeConfigFileNames: getDefaultBiomeConfigFiles(),
	}
}

// LoadBiome will search for the biome in the files specified by
//    BiomeParser.BiomeConfigFileNames
func (p *BiomeParser) LoadBiome(biomeName string) error {
	p.LoadedBiome = nil // Make sure we start from scratch

	for _, fPath := range p.BiomeConfigFileNames {
		// Does this file exist?
		if !fileExists(fPath) {
			continue
		}

		// Try to find the biome in that file
		freader, err := os.Open(fPath)
		if err != nil {
			continue
		}

		if p.tryLoadBiomeFromFile(biomeName, freader) {
			return nil
		}
	}

	return fmt.Errorf("unable to find the '%s' biome", biomeName)
}

// loadBiomeFile will load in a biome
func (p *BiomeParser) tryLoadBiomeFromFile(biome string, fcontents io.Reader) bool {
	buf := new(bytes.Buffer)
	buf.ReadFrom(fcontents)

	foundBiome, err := findBiomeInFileContents(buf.Bytes(), biome)
	if err != nil || foundBiome == nil {
		return false
	}

	p.LoadedBiome = foundBiome

	return true
}

// ConfigureBiome will do any biome configuration needed from the loaded biome
func (p BiomeParser) ConfigureBiome() error {
	// Was the biome loaded in?
	if p.LoadedBiome == nil {
		return fmt.Errorf("biome not loaded")
	}

	// Does an AWS session need to be configured?
	if p.LoadedBiome.AwsProfile != "" {
		err := config.ConfigureAwsEnvironment(p.LoadedBiome.AwsProfile)
		if err != nil {
			return fmt.Errorf("unable to configure AWS environment '%s': %v", p.LoadedBiome.AwsProfile, err)
		}
	}

	// Get the environment variables
	envs := p.LoadedBiome.GetEnvs()

	// Set the environment variables
	for key, val := range envs {
		os.Setenv(key, val)
	}

	return nil
}

// HELPERS

func findBiomeInFileContents(data []byte, biomeName string) (*config.BiomeConfig, error) {

	// Setup the reader and decoder
	r := bytes.NewReader(data)
	decoder := yaml.NewDecoder(r)

	// Loop over the yaml documents in the file
	for {
		var biome config.BiomeConfig

		if err := decoder.Decode(&biome); err != nil {
			if err != io.EOF {
				return nil, err
			}

			break
		}

		if biome.Name == biomeName {
			return &biome, nil
		}
	}

	return nil, nil
}

func getDefaultBiomeConfigFiles() []string {
	return []string{
		getCdFilePath(biomeConfigDefaultFileNames[0]),
		getCdFilePath(biomeConfigDefaultFileNames[1]),
		getHomeDirFilePath(biomeConfigDefaultFileNames[0]),
		getHomeDirFilePath(biomeConfigDefaultFileNames[1]),
	}
}

func getCdFilePath(fname string) string {
	cdPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get current directory: %v", err)
	}

	return path.Join(cdPath, fname)
}

func getHomeDirFilePath(fname string) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("unable to get current user: %v", err)
	}

	return path.Join(currentUser.HomeDir, fname)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}
