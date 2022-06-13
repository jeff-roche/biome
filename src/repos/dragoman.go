package repos

import "github.com/meltwater/dragoman/cryptography"

// The interface for the Dragoman Repository
type DragomanRepoIfc interface {
	Decrypt(string) (string, error)
}

// DragomanRepo handles decryption via dragoman
type DragomanRepo struct {
	client interface{}
}

// NewDragomanRepo is the builder function for DragomanRepo
func NewDragomanRepo() (*DragomanRepo, error) {
	return &DragomanRepo{
		client: nil,
	}, nil
}

func (r DragomanRepo) Decrypt(val string) (string, error) {

	// Setup the allowed cryptography techniques
	// Currently allow KMS Envelope encryption and Secrets Manager
	strat, err := cryptography.NewWildcardDecryptionStrategy([]cryptography.StrategyBuilder{
		func() (cryptography.Decryptor, error) { return cryptography.NewKmsCryptoStrategy("") },
		func() (cryptography.Decryptor, error) { return cryptography.NewSecretsManagerCryptoStrategy("") },
	})

	if err != nil {
		return "", err
	}

	data, err := strat.Decrypt(val)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
