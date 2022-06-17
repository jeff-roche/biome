package repos

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/jeff-roche/biome/src/lib/fileio"
	"github.com/jeff-roche/biome/src/lib/types"
	"gopkg.in/yaml.v3"
)

type BiomeFileParserIfc interface {
	FindBiome(biomeName string, searchFiles []string) (*types.BiomeConfig, error)
}

type BiomeFileParser struct{}

func NewBiomeFileParser() *BiomeFileParser {
	return &BiomeFileParser{}
}

func (parser BiomeFileParser) FindBiome(biomeName string, searchFiles []string) (*types.BiomeConfig, error) {
	// If we have any errors with the file, continue on to the next one
	for _, fPath := range searchFiles {
		if !fileio.FileExists(fPath) {
			continue
		}

		// Open the file and try to load the biome from it
		freader, err := os.Open(fPath)
		if err != nil {
			continue
		}

		biomes := parser.loadBiomes(freader)
		if _, exists := biomes[biomeName]; exists {
			biome := biomes[biomeName]

			if err := biome.Inherit(biomes); err != nil {
				return nil, err
			}

			return biome, nil
		}
	}

	return nil, fmt.Errorf("unable to locate the '%s' biome", biomeName)
}

// Load biomes will load in all biomes from the given io.Reader
func (parser BiomeFileParser) loadBiomes(fcontents io.Reader) map[string]*types.BiomeConfig {
	buff := new(bytes.Buffer)
	buff.ReadFrom(fcontents)

	// Loop over the documents in the .biome.yaml config
	reader := bytes.NewReader(buff.Bytes())
	decoder := yaml.NewDecoder(reader)

	biomes := make(map[string]*types.BiomeConfig)

	for {
		var biomeCfg types.BiomeConfig

		if err := decoder.Decode(&biomeCfg); err != nil {
			if err != io.EOF {
				continue
			}

			break
		}

		if biomeCfg.Name == "" {
			continue
		}

		biomes[biomeCfg.Name] = &biomeCfg
	}

	return biomes
}

// loadBiomeFromFile will search for the biome in the file and if it finds it will parse and return it
func (parser BiomeFileParser) loadBiomeFromFile(biomeName string, fcontents io.Reader) *types.BiomeConfig {
	buff := new(bytes.Buffer)
	buff.ReadFrom(fcontents)

	// Loop over the documents in the .biome.yaml config
	reader := bytes.NewReader(buff.Bytes())
	decoder := yaml.NewDecoder(reader)

	for {
		var biomeCfg types.BiomeConfig

		if err := decoder.Decode(&biomeCfg); err != nil {
			if err != io.EOF {
				continue
			}

			break
		}

		if biomeCfg.Name == biomeName {
			return &biomeCfg
		}
	}

	return nil
}
