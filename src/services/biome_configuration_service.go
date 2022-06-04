package services

import (
	"fmt"
	"path"

	"github.com/jeff-roche/biome/src/lib/fileio"
	"github.com/jeff-roche/biome/src/lib/setters"
	"github.com/jeff-roche/biome/src/lib/types"
	"github.com/jeff-roche/biome/src/repos"
)

var defaultFileNames = []string{".biome.yaml", ".biome.yml"}

// BiomeConfigurationService handles the loading and activation of biomes
type BiomeConfigurationService struct {
	ActiveBiome    *types.Biome
	configFileRepo repos.BiomeFileParserIfc
	awsStsRepo     repos.AwsStsRepositoryIfc
}

// NewBiomeConfigurationService is a builder function to generate the service
func NewBiomeConfigurationService() *BiomeConfigurationService {
	return &BiomeConfigurationService{
		configFileRepo: repos.NewBiomeFileParser(),
		awsStsRepo:     repos.NewAwsStsRepository(),
	}
}

// LoadBiome will search for the biome in the default locations
//     - Current directory .biome.[yaml|yml]
//     - Current user's home directory .biome.[yaml|yml]
func (svc *BiomeConfigurationService) LoadBiomeFromDefaults(biomeName string) error {
	// Setup the valid paths
	var validPaths []string
	if dir, err := fileio.GetCD(); err == nil {
		for _, fname := range defaultFileNames {
			validPaths = append(validPaths, path.Join(dir, fname))
		}
	}

	if dir, err := fileio.GetHomeDir(); err == nil {
		for _, fname := range defaultFileNames {
			validPaths = append(validPaths, path.Join(dir, fname))
		}
	}

	// Start blasting
	biome, err := svc.configFileRepo.FindBiome(biomeName, validPaths)
	if err != nil {
		svc.ActiveBiome = nil
		return err
	}

	svc.ActiveBiome = biome

	return nil
}

// LoadBiomeFromFile will search for the biome in the file specified
func (svc *BiomeConfigurationService) LoadBiomeFromFile(biomeName string, fpath string) error {
	biome, err := svc.configFileRepo.FindBiome(biomeName, []string{fpath})
	if err != nil {
		svc.ActiveBiome = nil
		return err
	}

	svc.ActiveBiome = biome

	return nil
}

func (svc *BiomeConfigurationService) ActivateBiome() error {
	if svc.ActiveBiome == nil {
		return fmt.Errorf("no biome loaded")
	}

	// AWS Profile Configuration
	if svc.ActiveBiome.Config.AwsProfile != "" {
		envCfg, err := svc.awsStsRepo.ConfigureSession(svc.ActiveBiome.Config.AwsProfile)
		if err != nil {
			return err
		}

		svc.awsStsRepo.SetAwsEnvs(envCfg)
	}

	// Loop over the envs and set them
	for env, val := range svc.ActiveBiome.Config.Environment {
		setter, err := setters.GetEnvironmentSetter(env, val)
		if err != nil {
			return fmt.Errorf("error setting '%s': %v", env, err)
		}

		err = setter.SetEnv()
		if err != nil {
			return fmt.Errorf("error setting '%s': %v", env, err)
		}
	}

	return nil
}
