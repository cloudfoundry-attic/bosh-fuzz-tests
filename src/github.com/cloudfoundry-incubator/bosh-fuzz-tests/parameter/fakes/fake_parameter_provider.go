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
	Template          *FakeTemplate
	Compilation       *FakeCompilation
	Update            *FakeUpdate
	CloudProperties   *FakeCloudProperties
	FixedMigratedFrom *FakeFixedMigratedFrom
	Variables         *FakeVariables
}

func NewFakeParameterProvider(persistentDiskDef string, vmTypeDef string) *FakeParameterProvider {
	return &FakeParameterProvider{
		Stemcell:          NewFakeStemcell(),
		PersistentDisk:    NewFakePersistentDisk(persistentDiskDef),
		VmType:            NewFakeVmType(vmTypeDef),
		AvailabilityZone:  NewFakeAvailabilityZone(),
		Template:          NewFakeTemplate(),
		Compilation:       NewFakeCompilation(),
		Update:            NewFakeUpdate(),
		CloudProperties:   NewFakeCloudProperties(),
		FixedMigratedFrom: NewFakeFixedMigratedFrom(),
		Variables:         NewFakeVariables(),
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
	} else if name == "network" {
		return p.Network
	} else if name == "template" {
		return p.Template
	} else if name == "compilation" {
		return p.Compilation
	} else if name == "update" {
		return p.Update
	} else if name == "cloud_properties" {
		return p.CloudProperties
	} else if name == "fixed_migrated_from" {
		return p.FixedMigratedFrom
	} else if name == "variables" {
		return p.FixedMigratedFrom
	}

	return nil
}
