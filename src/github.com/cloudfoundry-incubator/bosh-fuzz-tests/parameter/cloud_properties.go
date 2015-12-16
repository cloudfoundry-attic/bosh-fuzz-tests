package parameter

import (
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"

	"math/rand"
)

type cloudProperties struct {
	numOfProperties []int
	nameGenerator   bftnamegen.NameGenerator
	reuseDecider    bftdecider.Decider
}

func NewCloudProperties(
	numOfProperties []int,
	nameGenerator bftnamegen.NameGenerator,
	reuseDecider bftdecider.Decider,
) Parameter {
	return &cloudProperties{
		numOfProperties: numOfProperties,
		nameGenerator:   nameGenerator,
		reuseDecider:    reuseDecider,
	}
}

func (c *cloudProperties) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for i, subject := range input.CloudConfig.AvailabilityZones {
		found, prevSubject := subject.FindIn(previousInput.CloudConfig.AvailabilityZones)
		input.CloudConfig.AvailabilityZones[i].CloudProperties = c.FuzzCloudProperties(found, prevSubject.CloudProperties)
	}

	for i, subject := range input.CloudConfig.VmTypes {
		found, prevSubject := subject.FindIn(previousInput.CloudConfig.VmTypes)
		input.CloudConfig.VmTypes[i].CloudProperties = c.FuzzCloudProperties(found, prevSubject.CloudProperties)
	}

	return input
}

func (c *cloudProperties) FuzzCloudProperties(foundPrevProperties bool, prevProperties map[string]string) map[string]string {
	if c.reuseDecider.IsYes() && foundPrevProperties {
		return prevProperties
	}

	properties := map[string]string{}
	currentNumOfProperties := c.numOfProperties[rand.Intn(len(c.numOfProperties))]
	for j := 0; j < currentNumOfProperties; j++ {
		key := c.nameGenerator.Generate(7)
		value := c.nameGenerator.Generate(7)
		properties[key] = value
	}

	return properties
}
