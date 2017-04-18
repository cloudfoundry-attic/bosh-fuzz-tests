package flow

import (
	"strings"
	"time"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type actionsFlow struct {
	flowNumber       int
	actionInfos      []ActionInfo
	actionFactory    bltaction.Factory
	cliRunnerFactory bltclirunner.Factory
}

func NewFlow(
	flowNumber int,
	actionInfos []ActionInfo,
	actionFactory bltaction.Factory,
	cliRunnerFactory bltclirunner.Factory,
) *actionsFlow {
	return &actionsFlow{
		flowNumber:       flowNumber,
		actionInfos:      actionInfos,
		actionFactory:    actionFactory,
		cliRunnerFactory: cliRunnerFactory,
	}
}

func (f *actionsFlow) Run(usingLegacyManifest bool) error {
	uuid, err := boshuuid.NewGenerator().Generate()
	if err != nil {
		return err
	}
	deploymentName := strings.Join([]string{"deployment", uuid}, "-")

	cliRunner := f.cliRunnerFactory.Create()

	for i, actionInfo := range f.actionInfos {
		action, err := f.actionFactory.Create(actionInfo.Name, f.flowNumber, deploymentName, cliRunner, usingLegacyManifest)
		if err != nil {
			return err
		}

		err = action.Execute()
		if err != nil {
			return err
		}

		if i < len(f.actionInfos)-1 {
			time.Sleep(time.Duration(actionInfo.DelayInMilliseconds) * time.Millisecond)
		}
	}

	return nil
}
