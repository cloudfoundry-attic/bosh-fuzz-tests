package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Analyzer interface {
	Analyze(inputs []bftinput.Input) []Case
}

type Case struct {
	Input              bftinput.Input
	Expectations       []bftexpectation.Expectation
	DeploymentWillFail bool
}

type analyzer struct {
	stemcellComparator       Comparator
	nothingChangedComparator Comparator
	variablesComparator      Comparator
}

func NewAnalyzer(logger boshlog.Logger) Analyzer {
	return &analyzer{
		stemcellComparator:       NewStemcellComparator(logger),
		nothingChangedComparator: NewNothingChangedComparator(),
		variablesComparator:      NewVariablesComparator(),
	}
}

func (a *analyzer) Analyze(inputs []bftinput.Input) []Case {
	cases := []Case{}

	for i := range inputs {
		expectations := []bftexpectation.Expectation{}
		deploymentWillFail := false

		if i != 0 {
			expectations = append(expectations, a.stemcellComparator.Compare(inputs[:i], inputs[i])...)
			expectations = append(expectations, a.nothingChangedComparator.Compare(inputs[:i], inputs[i])...)

			deploymentWillFail = a.isMigratingFromAzsToNoAzsAndReusingStaticIps(inputs[i-1], inputs[i])
			deploymentWillFail = deploymentWillFail || a.isMovingInstancesStaticIPToAnotherAZ(inputs[i-1], inputs[i])
		}
		expectations = append(expectations, a.variablesComparator.Compare(inputs[:i], inputs[i])...)

		cases = append(cases, Case{
			Input:              inputs[i],
			Expectations:       expectations,
			DeploymentWillFail: deploymentWillFail,
		})
	}

	return cases
}

func (a *analyzer) isMigratingFromAzsToNoAzsAndReusingStaticIps(previousInput bftinput.Input, currentInput bftinput.Input) bool {
	for _, job := range currentInput.Jobs {
		previousJob, found := previousInput.FindJobByName(job.Name)
		if found && (len(previousJob.AvailabilityZones) > 0 && len(job.AvailabilityZones) == 0) {
			for _, network := range job.Networks {
				previousNetwork, networkFound := previousJob.FindNetworkByName(network.Name)
				if networkFound {
					for _, currentIP := range network.StaticIps {
						for _, prevIP := range previousNetwork.StaticIps {
							if prevIP == currentIP {
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

func (a *analyzer) isMovingInstancesStaticIPToAnotherAZ(previousInput bftinput.Input, currentInput bftinput.Input) bool {
	var previouslyUsedAZs []string

	for _, job := range previousInput.Jobs {
		for _, az := range job.AvailabilityZones {
			previouslyUsedAZs = append(previouslyUsedAZs, az)
		}
	}

	var currentlyAllowedAZs []string

	for _, az := range currentInput.CloudConfig.AvailabilityZones {
		currentlyAllowedAZs = append(currentlyAllowedAZs, az.Name)
	}

	var missingAZs []string

	for _, previouslyUsedAZ := range previouslyUsedAZs {
		previouslyUsedAZMissing := true

		for _, allowedAz := range currentlyAllowedAZs {
			if previouslyUsedAZ == allowedAz {
				previouslyUsedAZMissing = false
			}
		}

		if previouslyUsedAZMissing {
			missingAZs = append(missingAZs, previouslyUsedAZ)
		}
	}

	if len(missingAZs) > 0 {
		return true
	}

	return false
}
