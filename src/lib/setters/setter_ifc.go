package setters

// EnvironmentSetter is the minimum contract needed to set an environment variable of any kind
type EnvironmentSetter interface {
	SetEnv() error
}
