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
	FindBiome(biomeName string, searchFiles []string) (*types.Biome, error)
}

type BiomeFileParser struct{}

func NewBiomeFileParser() *BiomeFileParser {
	return &BiomeFileParser{}
}

func (parser BiomeFileParser) FindBiome(biomeName string, searchFiles []string) (*types.Biome, error) {
	var biome types.Biome // Setup a new biome

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

		biomeConfig := parser.loadBiomeFromFile(biomeName, freader)
		if biomeConfig != nil {
			biome.Config = biomeConfig
			biome.SourceFile = fPath

			return &biome, nil
		}
	}

	return nil, fmt.Errorf("unable to locate the '%s' biome", biomeName)
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
