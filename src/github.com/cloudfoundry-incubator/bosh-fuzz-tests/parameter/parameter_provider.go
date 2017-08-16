package parameter

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	bftnetwork "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type ParameterProvider interface {
	Get(name string) Parameter
}

type parameterProvider struct {
	parameters      bftconfig.Parameters
	nameGenerator   bftnamegen.NameGenerator
	reuseDecider    bftdecider.Decider
	networkAssigner bftnetwork.Assigner
	logger          boshlog.Logger
}

func NewParameterProvider(
	parameters bftconfig.Parameters,
	nameGenerator bftnamegen.NameGenerator,
	reuseDecider bftdecider.Decider,
	networkAssigner bftnetwork.Assigner,
	logger boshlog.Logger,
) ParameterProvider {
	return &parameterProvider{
		parameters:      parameters,
		nameGenerator:   nameGenerator,
		reuseDecider:    reuseDecider,
		networkAssigner: networkAssigner,
		logger:          logger,
	}
}

func (p *parameterProvider) Get(name string) Parameter {
	switch name {
	case "stemcell":
		stemcellDefinition := p.parameters.StemcellDefinition[rand.Intn(len(p.parameters.StemcellDefinition))]
		return NewStemcell(stemcellDefinition, p.parameters.StemcellVersions)
	case "persistent_disk":
		persistentDiskDefinition := p.parameters.PersistentDiskDefinition[rand.Intn(len(p.parameters.StemcellDefinition))]
		return NewPersistentDisk(persistentDiskDefinition, p.parameters.PersistentDiskSize, p.nameGenerator)
	case "vm_type":
		vmTypeDefinition := p.parameters.VmTypeDefinition[rand.Intn(len(p.parameters.VmTypeDefinition))]
		return NewVmType(vmTypeDefinition, p.nameGenerator, p.reuseDecider, p.logger)
	case "availability_zone":
		return NewAvailabilityZone(p.parameters.AvailabilityZones)
	case "network":
		return NewNetwork(p.networkAssigner)
	case "job":
		return NewJob(p.parameters.Jobs)
	case "compilation":
		return NewCompilation(p.parameters.NumberOfCompilationWorkers)
	case "update":
		return NewUpdate(p.parameters.Canaries, p.parameters.MaxInFlight, p.parameters.Serial)
	case "cloud_properties":
		return NewCloudProperties(p.parameters.NumOfCloudProperties, p.nameGenerator, p.reuseDecider)
	case "fixed_migrated_from":
		return NewFixedMigratedFrom()
	case "variables":
		numOfVariables := p.parameters.NumOfVariables[rand.Intn(len(p.parameters.NumOfVariables))]
		return NewVariables(numOfVariables, p.parameters.VariableTypes, p.nameGenerator, p.reuseDecider)
	case "lifecycle":
		return NewLifecycle()
	}

	return nil
}
