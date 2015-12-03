package parameter

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type ParameterProvider interface {
	Get(name string) Parameter
}

type parameterProvider struct {
	parameters    bftconfig.Parameters
	nameGenerator bftnamegen.NameGenerator
	reuseDecider  bftdecider.Decider
	logger        boshlog.Logger
}

func NewParameterProvider(
	parameters bftconfig.Parameters,
	nameGenerator bftnamegen.NameGenerator,
	reuseDecider bftdecider.Decider,
	logger boshlog.Logger,
) ParameterProvider {
	return &parameterProvider{
		parameters:    parameters,
		nameGenerator: nameGenerator,
		reuseDecider:  reuseDecider,
		logger:        logger,
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
		return NewVmType(vmTypeDefinition, p.nameGenerator, p.reuseDecider, p.logger)
	} else if name == "availability_zone" {
		return NewAvailabilityZone(p.parameters.AvailabilityZones)
	}

	return nil
}
