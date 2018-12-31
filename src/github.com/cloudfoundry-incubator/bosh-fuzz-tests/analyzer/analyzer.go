package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Analyzer interface {
	Analyze(inputs []bftinput.Input) []Case
}

type Case struct {
	Input                bftinput.Input
	Expectations         []bftexpectation.Expectation
	DeploymentWillFail   bool
	InstancesAfterDeploy map[string][]bltaction.Instance
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

	for i, input := range inputs {
		expectations := []bftexpectation.Expectation{}
		deploymentWillFail := false

		if !input.IsDryRun {
			if i != 0 {
				filteredInputs := filterDryRun(inputs[:i])
				expectations = append(expectations, a.stemcellComparator.Compare(filteredInputs, input)...)
				expectations = append(expectations, a.nothingChangedComparator.Compare(filteredInputs, input)...)

				filteredLastInput := filteredInputs[len(filteredInputs)-1]
				deploymentWillFail = a.isReusingStaticIps(filteredLastInput, input)
			}
			expectations = append(expectations, a.variablesComparator.Compare(nil, input)...)
			deploymentWillFail = deploymentWillFail || a.hasVariablesCertificateWithoutCA(input)
		}

		cases = append(cases, Case{
			Input:              input,
			Expectations:       expectations,
			DeploymentWillFail: deploymentWillFail,
		})
	}

	return cases
}

func filterDryRun(inputs []bftinput.Input) []bftinput.Input {
	output := []bftinput.Input{}
	for _, input := range inputs {
		if !input.IsDryRun {
			output = append(output, input)
		}
	}
	return output
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

func (a *analyzer) isReusingStaticIps(previousInput bftinput.Input, currentInput bftinput.Input) bool {
	return a.isMigratingFromAzsToNoAzsAndReusingStaticIps(previousInput, currentInput) ||
		a.isMovingInstancesStaticIPToAnotherAZ(previousInput, currentInput) ||
		a.isMovingInstancesStaticIPToAnotherInstanceGroup(previousInput, currentInput)
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

func (a *analyzer) isMovingInstancesStaticIPToAnotherInstanceGroup(previousInput bftinput.Input, currentInput bftinput.Input) bool {
	for _, instanceGroup := range currentInput.InstanceGroups {
		for _, previousInstanceGroup := range previousInput.InstanceGroups {
			if previousInstanceGroup.Name == instanceGroup.Name {
				continue
			}

			for _, network := range instanceGroup.Networks {
				for _, currentIP := range network.StaticIps {
					for _, previousNetwork := range previousInstanceGroup.Networks {
						for _, previousIP := range previousNetwork.StaticIps {
							if currentIP == previousIP {
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
