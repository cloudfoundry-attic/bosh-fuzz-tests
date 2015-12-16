package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakePersistentDisk struct {
}

func NewFakePersistentDisk() *FakePersistentDisk {
	return &FakePersistentDisk{}
}

func (s *FakePersistentDisk) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	persistentDiskPool := bftinput.DiskConfig{Name: "fake-persistent-disk", Size: 1}

	input.CloudConfig.PersistentDiskPools = []bftinput.DiskConfig{
		persistentDiskPool,
	}
	for j, _ := range input.Jobs {
		input.Jobs[j].PersistentDiskPool = persistentDiskPool.Name
	}

	return input
}
