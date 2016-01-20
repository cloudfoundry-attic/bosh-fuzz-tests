package flow

import (
	"encoding/json"
	"math/rand"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type ActionInfo struct {
	Name                string `json:"name"`
	DelayInMilliseconds int64  `json:"delay"`
}

type randomizer struct {
	actionFactory          bltaction.Factory
	cliRunnerFactory       bltclirunner.Factory
	state                  [][]ActionInfo
	maxDelayInMilliseconds int64
	fs                     boshsys.FileSystem
	logger                 boshlog.Logger
}

type Randomizer interface {
	Configure(filePath string) error
	Prepare(flows [][]string, numberOfDeployments int) error
	RunFlow(flowNumber int, usingLegacyManifest bool) error
}

func NewRandomizer(
	actionFactory bltaction.Factory,
	cliRunnerFactory bltclirunner.Factory,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) Randomizer {
	return &randomizer{
		actionFactory:    actionFactory,
		cliRunnerFactory: cliRunnerFactory,
		state:            [][]ActionInfo{},
		maxDelayInMilliseconds: 5000,
		fs:     fs,
		logger: logger,
	}
}

func (r *randomizer) Configure(filePath string) error {
	stateJSON, err := r.fs.ReadFile(filePath)
	if err != nil {
		return err
	}

	r.logger.Debug("randomizer", "Using pre-loaded state '%s'", stateJSON)

	err = json.Unmarshal([]byte(stateJSON), &r.state)
	if err != nil {
		return err
	}

	return nil
}

func (r *randomizer) Prepare(flows [][]string, numberOfDeployments int) error {
	for i := 0; i < numberOfDeployments; i++ {
		actionInfos := []ActionInfo{}

		randomActionNames := flows[rand.Intn(len(flows)-1)]

		for _, actionName := range randomActionNames {
			actionInfos = append(actionInfos, ActionInfo{
				Name:                actionName,
				DelayInMilliseconds: rand.Int63n(r.maxDelayInMilliseconds),
			})
		}
		r.state = append(r.state, actionInfos)
	}

	stateJSON, err := json.Marshal(r.state)
	if err != nil {
		return err
	}

	r.logger.Debug("randomizer", "Generated state '%s'", stateJSON)

	return nil
}

func (r *randomizer) RunFlow(flowNumber int, usingLegacyManifest bool) error {
	actionNames := r.state[flowNumber]
	r.logger.Debug("randomizer", "Creating flow with %#v", actionNames)

	flow := NewFlow(flowNumber, actionNames, r.actionFactory, r.cliRunnerFactory)

	return flow.Run(usingLegacyManifest)
}
