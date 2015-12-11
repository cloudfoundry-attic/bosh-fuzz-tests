package network

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type reuser struct {
	previousNetworks []bftinput.NetworkConfig
	decider          bftdecider.Decider
	nameGenerator    bftnamegen.NameGenerator
}

func NewReuser(
	previousNetworks []bftinput.NetworkConfig,
	decider bftdecider.Decider,
	nameGenerator bftnamegen.NameGenerator,
) *reuser {
	return &reuser{
		previousNetworks: previousNetworks,
		decider:          decider,
		nameGenerator:    nameGenerator,
	}
}

func (r *reuser) CreateNetwork(networkType string) bftinput.NetworkConfig {
	var network bftinput.NetworkConfig

	reusePreviousNetwork := r.decider.IsYes()
	if reusePreviousNetwork && len(r.previousNetworks) > 0 {
		networkToReuseIdx := rand.Intn(len(r.previousNetworks))
		network = bftinput.NetworkConfig{
			Name: r.previousNetworks[networkToReuseIdx].Name,
		}
		r.previousNetworks = append(r.previousNetworks[:networkToReuseIdx], r.previousNetworks[networkToReuseIdx+1:]...)

	} else {
		network = bftinput.NetworkConfig{
			Name: r.nameGenerator.Generate(7),
		}
	}

	network.Type = networkType
	return network
}
