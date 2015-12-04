package expectation

type Expectation interface {
	Run(taskId string) error
}
