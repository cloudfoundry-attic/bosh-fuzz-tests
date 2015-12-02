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

func (s *persistentDisk) Apply(input *bftinput.Input) *bftinput.Input {
	persistentDiskSize := s.diskSizes[rand.Intn(len(s.diskSizes))]

	for j, _ := range input.Jobs {
		if persistentDiskSize != 0 {
			if s.definition == "disk_pool" {
				input.Jobs[j].PersistentDiskPool = s.nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskPools = append(
					input.CloudConfig.PersistentDiskPools,
					bftinput.DiskConfig{Name: input.Jobs[j].PersistentDiskPool, Size: persistentDiskSize},
				)
			} else if s.definition == "disk_type" {
				input.Jobs[j].PersistentDiskType = s.nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskTypes = append(
					input.CloudConfig.PersistentDiskTypes,
					bftinput.DiskConfig{Name: input.Jobs[j].PersistentDiskType, Size: persistentDiskSize},
				)
			} else {
				input.Jobs[j].PersistentDiskSize = persistentDiskSize
			}
		}
	}

	return input
}
