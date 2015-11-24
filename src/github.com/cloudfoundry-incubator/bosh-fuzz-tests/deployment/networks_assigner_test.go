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
		networks = [][]string{[]string{"manual"}}
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
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
							},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					Networks: []NetworkConfig{
						{
							Name: "foo-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange:           "192.168.0.0/24",
									Gateway:           "192.168.0.1",
									AvailabilityZones: []string{"z1"},
									Reserved: []string{
										"192.168.0.71-192.168.0.75",
										"192.168.0.78-192.168.0.88",
										"192.168.0.90-192.168.0.113",
										"192.168.0.190",
										"192.168.0.231",
									},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange: "192.168.1.0/24",
									Gateway: "192.168.1.254",
									Reserved: []string{
										"192.168.1.153",
									},
								},
							},
						},
					},
					CompilationNetwork: "bar-net",
				},
			},
		},
		))
	})

	It("generates new subnet range for each subnet", func() {
		inputs := []Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         1,
						AvailabilityZones: []string{"z1"},
					},
					{
						Name:              "bar",
						Instances:         1,
						AvailabilityZones: []string{"z2"},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
				},
			},
		}
		networksAssigner.Assign(inputs)

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         1,
						AvailabilityZones: []string{"z1"},
						Networks: []JobNetworkConfig{
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
							},
						},
					},
					{
						Name:              "bar",
						Instances:         1,
						AvailabilityZones: []string{"z2"},
						Networks: []JobNetworkConfig{
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
							},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
					Networks: []NetworkConfig{
						{
							Name: "foo-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange:           "192.168.0.0/24",
									Gateway:           "192.168.0.1",
									AvailabilityZones: []string{"z2"},
									Reserved: []string{
										"192.168.0.85-192.168.0.98",
										"192.168.0.100-192.168.0.133",
									},
								},
								{
									IpRange:           "192.168.1.0/24",
									Gateway:           "192.168.1.254",
									AvailabilityZones: []string{"z2", "z1"},
									Reserved: []string{
										"192.168.1.132-192.168.1.142",
										"192.168.1.144-192.168.1.161",
										"192.168.1.170",
									},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange:  "192.168.2.0/24",
									Gateway:  "192.168.2.1",
									Reserved: []string{"192.168.2.243"},
								},
								{
									IpRange: "192.168.3.0/24",
									Gateway: "192.168.3.254",
									Reserved: []string{
										"192.168.3.224-192.168.3.227",
										"192.168.3.229-192.168.3.242",
										"192.168.3.246",
										"192.168.3.251",
									},
								},
							},
						},
					},
					CompilationNetwork:          "foo-net",
					CompilationAvailabilityZone: "z1",
				},
			},
		},
		))
	})
})
