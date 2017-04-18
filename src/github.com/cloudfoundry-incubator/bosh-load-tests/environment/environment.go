package environment

type Environment interface {
	Setup() error
	Shutdown() error
	DirectorURL() string
}
