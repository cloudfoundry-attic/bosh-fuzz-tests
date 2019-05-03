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
		prevSubject, found := previousInput.FindAzByName(subject.Name)
		input.CloudConfig.AvailabilityZones[i].CloudProperties = c.FuzzCloudProperties(found, prevSubject.CloudProperties)
	}

	for i, subject := range input.CloudConfig.VmTypes {
		prevSubject, found := previousInput.FindVmTypeByName(subject.Name)
		input.CloudConfig.VmTypes[i].CloudProperties = c.FuzzCloudProperties(found, prevSubject.CloudProperties)
	}

	for i, subject := range input.CloudConfig.PersistentDiskTypes {
		prevSubject, found := previousInput.FindDiskTypeByName(subject.Name)
		input.CloudConfig.PersistentDiskTypes[i].CloudProperties = c.FuzzCloudProperties(found, prevSubject.CloudProperties)
	}

	// we can't really detect when a previous input has used cloud properties since we could
	// validly used 0 properties
	input.CloudConfig.Compilation.CloudProperties = c.FuzzCloudProperties(true, previousInput.CloudConfig.Compilation.CloudProperties)

	for i, network := range input.CloudConfig.Networks {
		if network.Subnets != nil {
			for s, subject := range network.Subnets {
				if subject.IpPool != nil {
					// manual
					prevSubject, found := previousInput.FindSubnetByIpRange(subject.IpPool.IpRange)
					input.CloudConfig.Networks[i].Subnets[s].CloudProperties = c.FuzzCloudProperties(found, prevSubject.CloudProperties)
				} else {
					// dynamic
					input.CloudConfig.Networks[i].Subnets[s].CloudProperties = c.FuzzCloudProperties(false, map[string]string{})
				}
			}
		} else {
			// vip
			input.CloudConfig.Networks[i].CloudProperties = c.FuzzCloudProperties(false, map[string]string{})
		}
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
