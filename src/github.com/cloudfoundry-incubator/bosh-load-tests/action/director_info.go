package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type DirectorInfo struct {
	UUID string
	URL  string
}

func NewDirectorInfo(directorURL string, cliRunnerFactory bltclirunner.Factory) (DirectorInfo, error) {
	cliRunner := cliRunnerFactory.Create()
	cliRunner.Configure()
	defer cliRunner.Clean()

	err := cliRunner.TargetAndLogin(directorURL)
	if err != nil {
		return DirectorInfo{}, err
	}

	uuid, err := cliRunner.RunWithOutput("status", "--uuid")
	if err != nil {
		return DirectorInfo{}, err
	}

	return DirectorInfo{
		UUID: uuid,
		URL:  directorURL,
	}, nil
}
