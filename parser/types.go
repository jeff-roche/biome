package parser

import "fmt"

type BiomeNotFoundError struct {
	Name string
}

func NewBiomeNotFoundError(biomeName string) BiomeNotFoundError {
	return BiomeNotFoundError{
		Name: biomeName,
	}
}

func (e BiomeNotFoundError) Error() string {
	return fmt.Sprintf("unable to find Biome '%s'", e.Name)
}