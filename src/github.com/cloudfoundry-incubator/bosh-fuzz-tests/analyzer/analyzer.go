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
		expectations = append(expectations, a.variablesComparator.Compare(nil, inputs[i])...)
		deploymentWillFail = deploymentWillFail || a.hasVariablesCertificateWithoutCA(inputs[i])

		cases = append(cases, Case{
			Input:              inputs[i],
			Expectations:       expectations,
			DeploymentWillFail: deploymentWillFail,
		})
	}

	return cases
}
func (a *analyzer) hasVariablesCertificateWithoutCA(currentInput bftinput.Input) bool {
	for _, variable := range currentInput.Variables {
		if variable.Type != "certificate" {
			continue
		}

		options := variable.Options
		isCA, ok := options["is_ca"].(bool)
		if !ok {
			isCA = false
		}

		referencedCA, found := options["ca"]
		if isCA && !found {
			continue
		}

		if !isCA && !found {
			return true
		}

		var referencedVariable bftinput.Variable
		for _, variable := range currentInput.Variables {
			if variable.Name == referencedCA {
				referencedVariable = variable
				break
			}
		}

		if referencedVariable.Type != "certificate" {
			return true
		} else {
			isCA, ok := referencedVariable.Options["is_ca"].(bool)
			if !ok || !isCA {
				return true
			}
		}
	}

	return false
}

func (a *analyzer) isMigratingFromAzsToNoAzsAndReusingStaticIps(previousInput bftinput.Input, currentInput bftinput.Input) bool {
	for _, instanceGroup := range currentInput.InstanceGroups {
		previousInstanceGroup, found := previousInput.FindInstanceGroupByName(instanceGroup.Name)
		if found && (len(previousInstanceGroup.AvailabilityZones) > 0 && len(instanceGroup.AvailabilityZones) == 0) {
			for _, network := range instanceGroup.Networks {
				previousNetwork, networkFound := previousInstanceGroup.FindNetworkByName(network.Name)
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

	for _, instanceGroup := range previousInput.InstanceGroups {
		for _, az := range instanceGroup.AvailabilityZones {
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
