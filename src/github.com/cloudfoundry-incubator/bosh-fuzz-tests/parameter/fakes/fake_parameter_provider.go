package fakes

import (
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
)

type FakeParameterProvider struct {
	Stemcell         *FakeStemcell
	PersistentDisk   *FakePersistentDisk
	VmType           *FakeVmType
	AvailabilityZone *FakeAvailabilityZone
}

func NewFakeParameterProvider() *FakeParameterProvider {
	return &FakeParameterProvider{
		Stemcell:         NewFakeStemcell(),
		PersistentDisk:   NewFakePersistentDisk(),
		VmType:           NewFakeVmType(),
		AvailabilityZone: NewFakeAvailabilityZone(),
	}
}

func (p *FakeParameterProvider) Get(name string) bftparam.Parameter {
	if name == "stemcell" {
		return p.Stemcell
	} else if name == "persistent_disk" {
		return p.PersistentDisk
	} else if name == "vm_type" {
		return p.VmType
	} else if name == "availability_zone" {
		return p.AvailabilityZone
	}

	return nil
}
