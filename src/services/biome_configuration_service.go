package services

import (
	"fmt"
	"path"
	"strings"

	"github.com/jeff-roche/biome/src/lib/cmdr"
	"github.com/jeff-roche/biome/src/lib/fileio"
	"github.com/jeff-roche/biome/src/lib/setters"
	"github.com/jeff-roche/biome/src/lib/types"
	"github.com/jeff-roche/biome/src/repos"
	"github.com/joho/godotenv"
)

var defaultFileNames = []string{".biome.yaml", ".biome.yml"}

// BiomeConfigurationService handles the loading and activation of biomes
type BiomeConfigurationService struct {
	ActiveBiome    *types.BiomeConfig
	configFileRepo repos.BiomeFileParserIfc
	awsStsRepo     repos.AwsStsRepositoryIfc
	configuredEnvs map[string]string
}

// NewBiomeConfigurationService is a builder function to generate the service
func NewBiomeConfigurationService() *BiomeConfigurationService {
	return &BiomeConfigurationService{
		configFileRepo: repos.NewBiomeFileParser(),
		awsStsRepo:     repos.NewAwsStsRepository(),
		configuredEnvs: make(map[string]string),
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

// SaveBiomeToFile will export the loaded environment variables to the file specified
func (svc BiomeConfigurationService) SaveBiomeToFile(fpath string) error {
	return godotenv.Write(svc.configuredEnvs, fpath)
}

// Activate biome will load up the configuration and run any setup commands before running the specified program
func (svc *BiomeConfigurationService) ActivateBiome() error {
	if svc.ActiveBiome == nil {
		return fmt.Errorf("no biome loaded")
	}

	// AWS
	if err := svc.loadAws(); err != nil {
		return err
	}

	// Dot Env
	if err := svc.loadFromEnv(svc.ActiveBiome.ExternalEnvFile); err != nil {
		return err
	}

	// Parse all Envs
	if err := svc.loadEnvs(); err != nil {
		return err
	}

	// Additional Commands
	if err := svc.runSetupCommands(); err != nil {
		return err
	}

	return nil
}

// loadAws will load in the AWS profile if one was specified
func (svc *BiomeConfigurationService) loadAws() error {
	if svc.ActiveBiome.AwsProfile != "" {
		envCfg, err := svc.awsStsRepo.ConfigureSession(svc.ActiveBiome.AwsProfile)
		if err != nil {
			return err
		}

		svc.awsStsRepo.SetAwsEnvs(envCfg)
	}

	return nil
}

// loadFromEnv will load in addition environment variables from the ENV file
//     Any envs specified in the biome config will override vars specified in the dotenv
func (svc *BiomeConfigurationService) loadFromEnv(fname string) error {
	if fname != "" {
		loadedEnvs, err := godotenv.Read(fname)
		if err != nil {
			return err
		}

		for key, val := range loadedEnvs {

			// Only save the key if one wasn't specified in the biome config
			if _, exists := svc.ActiveBiome.Environment[key]; !exists {
				svc.ActiveBiome.Environment[key] = val
			}
		}
	}

	return nil
}

// loadEnvs will parse all the envs in the Environment map and load them into memory
func (svc *BiomeConfigurationService) loadEnvs() error {
	// Loop over the envs and set them
	for env, val := range svc.ActiveBiome.Environment {
		setter, err := setters.GetEnvironmentSetter(env, val)
		if err != nil {
			return fmt.Errorf("error setting '%s': %v", env, err)
		}

		raw_val, err := setter.SetEnv()
		if err != nil {
			return fmt.Errorf("error setting '%s': %v", env, err)
		}

		// Save off the envs we configured
		svc.configuredEnvs[env] = raw_val
	}

	return nil
}

// runSetupCommands will run any command line commands specified in the biome configuration
func (svc *BiomeConfigurationService) runSetupCommands() error {
	if len(svc.ActiveBiome.Commands) > 0 {
		for _, cmd := range svc.ActiveBiome.Commands {
			parts := strings.Split(cmd, " ")

			if err := cmdr.Run(parts[0], parts[1:]...); err != nil {
				return err
			}
		}
	}

	return nil
}
