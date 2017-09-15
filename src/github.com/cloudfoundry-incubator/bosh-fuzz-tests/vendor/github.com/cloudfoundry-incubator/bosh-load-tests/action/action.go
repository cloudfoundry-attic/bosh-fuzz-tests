package action

type Action interface {
	Execute() error
}
