package dummy

type Database interface {
	Create() error
	Drop() error
	Name() string
	Server() string
	User() string
	Password() string
	Port() int
}
