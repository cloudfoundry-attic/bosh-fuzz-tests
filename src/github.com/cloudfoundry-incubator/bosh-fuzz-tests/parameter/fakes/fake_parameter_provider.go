package fakes

import (
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
)

type FakeParameterProvider struct {
	Stemcell          *FakeStemcell
	PersistentDisk    *FakePersistentDisk
	VmType            *FakeVmType
	AvailabilityZone  *FakeAvailabilityZone
	Network           *FakeNetwork
	Job               *FakeJob
	Compilation       *FakeCompilation
	Update            *FakeUpdate
	CloudProperties   *FakeCloudProperties
	FixedMigratedFrom *FakeFixedMigratedFrom
	Variables         *FakeVariables
	Lifecycle         *FakeLifecycle
}

func NewFakeParameterProvider(persistentDiskDef string) *FakeParameterProvider {
	return &FakeParameterProvider{
		Stemcell:          NewFakeStemcell(),
		PersistentDisk:    NewFakePersistentDisk(persistentDiskDef),
		VmType:            NewFakeVmType(),
		AvailabilityZone:  NewFakeAvailabilityZone(),
		Job:               NewFakeJob(),
		Compilation:       NewFakeCompilation(),
		Update:            NewFakeUpdate(),
		CloudProperties:   NewFakeCloudProperties(),
		FixedMigratedFrom: NewFakeFixedMigratedFrom(),
		Variables:         NewFakeVariables(),
		Lifecycle:         NewFakeLifecycle(),
	}
}

func (p *FakeParameterProvider) Get(name string) bftparam.Parameter {
	switch name {
	case "stemcell":
		return p.Stemcell
	case "persistent_disk":
		return p.PersistentDisk
	case "vm_type":
		return p.VmType
	case "availability_zone":
		return p.AvailabilityZone
	case "network":
		return p.Network
	case "job":
		return p.Job
	case "compilation":
		return p.Compilation
	case "update":
		return p.Update
	case "cloud_properties":
		return p.CloudProperties
	case "fixed_migrated_from":
		return p.FixedMigratedFrom
	case "variables":
		return p.FixedMigratedFrom
	case "lifecycle":
		return p.Lifecycle
	}

	return nil
}
