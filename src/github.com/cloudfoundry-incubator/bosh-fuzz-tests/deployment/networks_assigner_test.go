package deployment_test

import (
	fakebftdepl "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworksAssigner", func() {
	var (
		networksAssigner NetworksAssigner
		networks         [][]string
	)

	BeforeEach(func() {
		networks = [][]string{[]string{"dynamic"}}
		nameGenerator := &fakebftdepl.FakeNameGenerator{}
		nameGenerator.Names = []string{"foo-net", "bar-net", "baz-net"}
		networksAssigner = NewSeededNetworksAssigner(networks, nameGenerator, 5)
	})

	It("assigns network of the given type to job and cloud config", func() {
		inputs := []Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
				},
			},
		}

		networksAssigner.Assign(inputs)

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []JobNetworkConfig{
							{Name: "foo-net"},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					Networks: []NetworkConfig{
						{
							Name: "foo-net",
							Type: "dynamic",
							Subnets: []SubnetConfig{
								{
									AvailabilityZones: []string{"z1"},
								},
							},
						},
						{
							Name: "default",
							Subnets: []SubnetConfig{
								{
									AvailabilityZones: []string{"z1"},
								},
							},
						},
					},
				},
			},
		},
		))

	})
})
