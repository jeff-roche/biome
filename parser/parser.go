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

var LoadedBiome *config.BiomeConfig
var BiomeConfigFileNames = []string{".biome.yaml", ".biome.yml"}

// LoadBiome will find the nearest .biome.[yml|yaml] in this priority order
//    - Current Directory
//    - Home Directory
func LoadBiome(biomeName string) error {
	biomeConfigPaths := []string{
		getCdFilePath(BiomeConfigFileNames[0]),
		getCdFilePath(BiomeConfigFileNames[1]),
		getHomeDirFilePath(BiomeConfigFileNames[0]),
		getHomeDirFilePath(BiomeConfigFileNames[1]),
	}

	// Reset the loaded biome
	LoadedBiome = nil

	for _, fPath := range biomeConfigPaths {
		err := tryLoadBiomeConfig(fPath, biomeName)
		if err != nil {
			continue
		}

		// Was the biome loaded?
		if LoadedBiome != nil {
			return nil
		}
	}

	return fmt.Errorf("no biome configuration found")
}

func ConfigureBiome() error {
	// Did the Biome get configured?
	if LoadedBiome == nil {
		return fmt.Errorf("biome not configured")
	}

	// Does an AWS session need to be configured?
	if LoadedBiome.AwsProfile != "" {
		err := config.ConfigureAwsEnvironment(LoadedBiome.AwsProfile)
		if err != nil {
			return fmt.Errorf("unable to configure AWS environment '%s': %v", LoadedBiome.AwsProfile, err)
		}
	}

	// Get the environment variables
	envs := LoadedBiome.GetEnvs()

	// Set the environment variables
	for key, val := range envs {
		os.Setenv(key, val)
	}

	return nil
}

func tryLoadBiomeConfig(fPath string, biomeName string) error {
	if !fileExists(fPath) {
		return fmt.Errorf("could not find file '%s'", fPath)
	}

	// Slurp slurp
	biomeConfigContents, err := os.ReadFile(fPath)
	if err != nil {
		return fmt.Errorf("unable to load the biome config '%s': %v", fPath, err)
	}

	// Loop over the file contents and try to parse out biomes
	biome, err := findBiomeInFileContents(biomeConfigContents, biomeName)
	if err != nil {
		return fmt.Errorf("error searching for biome in '%s': %v", fPath, err)
	}

	// Save off the biome if we found it
	if biome != nil {
		LoadedBiome = biome
	}

	return nil
}

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
