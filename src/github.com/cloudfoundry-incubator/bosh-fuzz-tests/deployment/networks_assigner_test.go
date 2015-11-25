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
		expectedIpPool   *IpPool
	)

	BeforeEach(func() {
		networks = [][]string{[]string{"manual", "vip"}}
		nameGenerator := &fakebftdepl.FakeNameGenerator{}
		nameGenerator.Names = []string{"foo-net", "bar-net", "baz-net", "qux-net"}
		ipPoolProvider := &fakebftdepl.FakeIpPoolProvider{}
		vipPool := &IpPool{
			AvailableIps: []string{
				"10.10.0.6",
				"10.10.0.32",
			},
		}
		ipPoolProvider.RegisterIpPool(vipPool)

		ipPool := &IpPool{
			IpRange: "192.168.0.0/24",
			Gateway: "192.168.0.1",
			Reserved: []string{
				"192.168.0.15-192.168.0.58",
				"192.168.0.157-192.168.0.203",
			},
			AvailableIps: []string{
				"192.168.0.222",
				"192.168.0.110",
			},
		}
		ipPoolProvider.RegisterIpPool(ipPool)
		ipPoolProvider.RegisterIpPool(ipPool)

		expectedIpPool = &IpPool{
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
			AvailableIps: []string{},
		}
		staticIpDecider := &fakebftdepl.FakeDecider{}
		staticIpDecider.IsYesYes = true
		networksAssigner = NewSeededNetworksAssigner(networks, nameGenerator, ipPoolProvider, staticIpDecider, 32)
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
								Name:      "bar-net",
								StaticIps: []string{"10.10.0.6", "10.10.0.32"},
							},
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
									IpPool:            expectedIpPool,
									AvailabilityZones: []string{"z1"},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "vip",
						},
						{
							Name: "baz-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpPool: expectedIpPool,
								},
							},
						},
						{
							Name: "qux-net",
							Type: "vip",
						},
					},
					CompilationNetwork: "baz-net",
				},
			},
		},
		))
	})
})
