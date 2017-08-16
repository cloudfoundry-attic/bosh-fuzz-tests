package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakePersistentDisk struct {
	definition string
}

func NewFakePersistentDisk(definition string) *FakePersistentDisk {
	return &FakePersistentDisk{
		definition: definition,
	}
}

func (s *FakePersistentDisk) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	if s.definition == "disk_pool" {
		persistentDiskPool := bftinput.DiskConfig{Name: "fake-persistent-disk", Size: 1}

		input.CloudConfig.PersistentDiskPools = []bftinput.DiskConfig{
			persistentDiskPool,
		}
		for j, _ := range input.InstanceGroups {
			input.InstanceGroups[j].PersistentDiskPool = persistentDiskPool.Name
		}
	} else if s.definition == "disk_type" {
		persistentDiskType := bftinput.DiskConfig{Name: "fake-persistent-disk", Size: 1}

		input.CloudConfig.PersistentDiskTypes = []bftinput.DiskConfig{
			persistentDiskType,
		}
		for j, _ := range input.InstanceGroups {
			input.InstanceGroups[j].PersistentDiskType = persistentDiskType.Name
		}
	} else {
		for j, _ := range input.InstanceGroups {
			input.InstanceGroups[j].PersistentDiskSize = 10
		}
	}

	return input
}
