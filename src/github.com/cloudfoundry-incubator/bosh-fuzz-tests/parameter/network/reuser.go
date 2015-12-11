package network

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type reuser struct {
	previousNetworkNames []string
	decider              bftdecider.Decider
	nameGenerator        bftnamegen.NameGenerator
}

func NewReuser(
	previousNetworks []bftinput.NetworkConfig,
	decider bftdecider.Decider,
	nameGenerator bftnamegen.NameGenerator,
) *reuser {
	previousNetworkNames := []string{}
	for _, networkNames := range previousNetworks {
		previousNetworkNames = append(previousNetworkNames, networkNames.Name)
	}
	return &reuser{
		previousNetworkNames: previousNetworkNames,
		decider:              decider,
		nameGenerator:        nameGenerator,
	}
}

func (r *reuser) CreateNetwork(networkType string) bftinput.NetworkConfig {
	var network bftinput.NetworkConfig

	reusePreviousNetwork := r.decider.IsYes()
	if reusePreviousNetwork && len(r.previousNetworkNames) > 0 {
		networkToReuseIdx := rand.Intn(len(r.previousNetworkNames))
		network = bftinput.NetworkConfig{
			Name: r.previousNetworkNames[networkToReuseIdx],
		}
		r.previousNetworkNames = append(r.previousNetworkNames[:networkToReuseIdx], r.previousNetworkNames[networkToReuseIdx+1:]...)

	} else {
		network = bftinput.NetworkConfig{
			Name: r.nameGenerator.Generate(7),
		}
	}

	network.Type = networkType
	return network
}
