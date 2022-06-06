package repos

import (
	"github.com/jeff-roche/biome/src/lib/types"
	"github.com/stretchr/testify/mock"
)

type MockBiomeFileParser struct {
	mock.Mock
}

func (m MockBiomeFileParser) FindBiome(biomeName string, searchFiles []string) (*types.Biome, error) {
	args := m.Called(biomeName, searchFiles)
	return args.Get(0).(*types.Biome), args.Error(1)
}

type MockAwsStsRepository struct {
	mock.Mock
}

func (m MockAwsStsRepository) ConfigureSession(profile string) (*types.AwsEnvConfig, error) {
	args := m.Called(profile)
	return args.Get(0).(*types.AwsEnvConfig), args.Error(1)
}

func (m MockAwsStsRepository) SetAwsEnvs(cfg *types.AwsEnvConfig) {}
