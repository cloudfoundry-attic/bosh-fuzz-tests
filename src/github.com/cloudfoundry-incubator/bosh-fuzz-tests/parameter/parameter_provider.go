package parameter

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type ParameterProvider interface {
	Get(name string) Parameter
}

type parameterProvider struct {
	parameters    bftconfig.Parameters
	nameGenerator bftnamegen.NameGenerator
}

func NewParameterProvider(parameters bftconfig.Parameters, nameGenerator bftnamegen.NameGenerator) ParameterProvider {
	return &parameterProvider{
		parameters:    parameters,
		nameGenerator: nameGenerator,
	}
}

func (p *parameterProvider) Get(name string) Parameter {
	if name == "stemcell" {
		stemcellDefinition := p.parameters.StemcellDefinition[rand.Intn(len(p.parameters.StemcellDefinition))]
		return NewStemcell(stemcellDefinition, p.parameters.StemcellVersions)
	} else if name == "persistent_disk" {
		persistentDiskDefinition := p.parameters.PersistentDiskDefinition[rand.Intn(len(p.parameters.StemcellDefinition))]
		return NewPersistentDisk(persistentDiskDefinition, p.parameters.PersistentDiskSize, p.nameGenerator)
	} else if name == "vm_type" {
		vmTypeDefinition := p.parameters.VmTypeDefinition[rand.Intn(len(p.parameters.VmTypeDefinition))]
		return NewVmType(vmTypeDefinition, p.nameGenerator)
	} else if name == "availability_zone" {
		return NewAvailabilityZone(p.parameters.AvailabilityZones)
	}

	return nil
}
