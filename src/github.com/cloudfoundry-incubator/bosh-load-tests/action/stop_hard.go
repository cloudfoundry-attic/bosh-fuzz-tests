package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type stopHard struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
}

func NewStopHard(directorInfo DirectorInfo, deploymentName string, cliRunner bltclirunner.Runner, fs boshsys.FileSystem) *stopHard {
	return &stopHard{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
	}
}

func (s *stopHard) Execute() error {
	s.cliRunner.SetEnv(s.directorInfo.URL)

	deployWrapper := NewDeployWrapper(s.cliRunner)
	_, err := deployWrapper.RunWithDebug("-d", s.deploymentName, "stop", "--hard", "simple/0")
	if err != nil {
		return err
	}

	return nil
}
