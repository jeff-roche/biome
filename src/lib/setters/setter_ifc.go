package setters

// EnvironmentSetter is the minimum contract needed to set an environment variable of any kind
type EnvironmentSetter interface {
	// Returns the value being set and an error if one was encountered
	SetEnv() (string, error)
}
