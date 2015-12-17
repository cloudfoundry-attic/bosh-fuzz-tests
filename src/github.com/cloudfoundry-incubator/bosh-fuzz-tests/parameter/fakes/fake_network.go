package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeNetwork struct{}

func NewFakeNetwork() *FakeNetwork {
	return &FakeNetwork{}
}

func (s *FakeNetwork) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.CloudConfig.Networks = []bftinput.NetworkConfig{
		{
			Name: "foo-network",
			Subnets: []bftinput.SubnetConfig{
				{
					IpPool: &bftinput.IpPool{
						IpRange: "10.0.0.0/24",
					},
				},
			},
		},
	}
	for j, _ := range input.Jobs {
		input.Jobs[j].Networks = []bftinput.JobNetworkConfig{
			{
				Name: "foo-network",
			},
		}
	}

	return input
}
