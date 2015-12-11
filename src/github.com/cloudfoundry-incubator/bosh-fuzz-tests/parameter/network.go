package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnetwork "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"
)

type network struct {
	networkAssigner bftnetwork.Assigner
}

func NewNetwork(
	networkAssigner bftnetwork.Assigner,
) Parameter {
	return &network{
		networkAssigner: networkAssigner,
	}
}

func (n *network) Apply(input bftinput.Input) bftinput.Input {
	return n.networkAssigner.Assign(input)
}
