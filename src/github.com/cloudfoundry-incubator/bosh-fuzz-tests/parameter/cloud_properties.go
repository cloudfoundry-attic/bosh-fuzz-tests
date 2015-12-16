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
	for i, az := range input.CloudConfig.AvailabilityZones {
		foundPrevProperties, prevProperties := c.FindCloudPropertiesByAzName(previousInput.CloudConfig.AvailabilityZones, az.Name)
		if c.reuseDecider.IsYes() && foundPrevProperties {
			input.CloudConfig.AvailabilityZones[i].CloudProperties = prevProperties
		} else {
			input.CloudConfig.AvailabilityZones[i].CloudProperties = map[string]string{}
			currentNumOfProperties := c.numOfProperties[rand.Intn(len(c.numOfProperties))]
			for j := 0; j < currentNumOfProperties; j++ {
				key := c.nameGenerator.Generate(7)
				value := c.nameGenerator.Generate(7)
				input.CloudConfig.AvailabilityZones[i].CloudProperties[key] = value
			}
		}
	}

	return input
}

func (c *cloudProperties) FindCloudPropertiesByAzName(previousAzs []bftinput.AvailabilityZone, azName string) (bool, map[string]string) {
	for _, prevAz := range previousAzs {
		if prevAz.Name == azName {
			return true, prevAz.CloudProperties
		}
	}

	return false, map[string]string{}
}
