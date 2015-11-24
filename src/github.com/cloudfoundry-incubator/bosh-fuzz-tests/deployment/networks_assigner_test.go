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

		Expect(inputs).To(BeEquivalentTo([]Input{
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
								StaticIps:     []string{"192.168.0.222", "192.168.0.110"},
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
									IpPool: &IpPool{
										IpRange: "192.168.0.0/24",
										Gateway: "192.168.0.1",
										Reserved: []string{
											"192.168.0.15-192.168.0.58",
											"192.168.0.157-192.168.0.203",
										},
										Static: []string{
											"192.168.0.222",
											"192.168.0.110",
										},
									},
									AvailabilityZones: []string{"z1"},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpPool: &IpPool{
										IpRange: "192.168.1.0/24",
										Gateway: "192.168.1.254",
										Reserved: []string{
											"192.168.1.41-192.168.1.111",
											"192.168.1.132",
											"192.168.1.235",
										},
										AvailableIps: []string{"192.168.1.154"},
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
									IpPool: &IpPool{
										IpRange: "192.168.0.0/24",
										Gateway: "192.168.0.1",
										Reserved: []string{
											"192.168.0.71-192.168.0.75",
											"192.168.0.78-192.168.0.88",
											"192.168.0.90-192.168.0.113",
											"192.168.0.190",
											"192.168.0.231",
										},
									},
									AvailabilityZones: []string{"z1"},
								},
								{
									IpPool: &IpPool{
										IpRange: "192.168.0.0/24",
										Gateway: "192.168.0.1",
										Reserved: []string{
											"192.168.0.71-192.168.0.75",
											"192.168.0.78-192.168.0.88",
											"192.168.0.90-192.168.0.113",
											"192.168.0.190",
											"192.168.0.231",
										},
									},
									AvailabilityZones: []string{"z1"},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpPool: &IpPool{
										IpRange: "192.168.0.0/24",
										Gateway: "192.168.0.1",
										Reserved: []string{
											"192.168.0.71-192.168.0.75",
											"192.168.0.78-192.168.0.88",
											"192.168.0.90-192.168.0.113",
											"192.168.0.190",
											"192.168.0.231",
										},
									},
									AvailabilityZones: []string{"z1"},
								},
								{
									IpPool: &IpPool{
										IpRange: "192.168.0.0/24",
										Gateway: "192.168.0.1",
										Reserved: []string{
											"192.168.0.71-192.168.0.75",
											"192.168.0.78-192.168.0.88",
											"192.168.0.90-192.168.0.113",
											"192.168.0.190",
											"192.168.0.231",
										},
									},
									AvailabilityZones: []string{"z1"},
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
