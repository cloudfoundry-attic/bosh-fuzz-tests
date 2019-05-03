package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type persistentDisk struct {
	definition    string
	diskSizes     []int
	nameGenerator bftnamegen.NameGenerator
}

func NewPersistentDisk(definition string, diskSizes []int, nameGenerator bftnamegen.NameGenerator) Parameter {
	return &persistentDisk{
		definition:    definition,
		diskSizes:     diskSizes,
		nameGenerator: nameGenerator,
	}
}

func (s *persistentDisk) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	persistentDiskSize := s.diskSizes[rand.Intn(len(s.diskSizes))]

	input.CloudConfig.PersistentDiskTypes = nil

	for j := range input.InstanceGroups {
		input.InstanceGroups[j].PersistentDiskSize = 0
		input.InstanceGroups[j].PersistentDiskType = ""

		if persistentDiskSize != 0 {
			if s.definition == "disk_type" {
				input.InstanceGroups[j].PersistentDiskType = s.nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskTypes = append(
					input.CloudConfig.PersistentDiskTypes,
					bftinput.DiskConfig{Name: input.InstanceGroups[j].PersistentDiskType, Size: persistentDiskSize},
				)
			} else {
				input.InstanceGroups[j].PersistentDiskSize = persistentDiskSize
			}
		}
	}

	return input
}
