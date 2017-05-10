package expectation

type Expectation interface {
	Run(debugLog string) error
}
