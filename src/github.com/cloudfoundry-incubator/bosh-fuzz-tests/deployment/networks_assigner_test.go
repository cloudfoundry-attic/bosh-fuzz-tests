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
		networks = [][]string{[]string{"manual"}}
		nameGenerator := &fakebftdepl.FakeNameGenerator{}
		nameGenerator.Names = []string{"foo-net", "bar-net", "baz-net"}
		ipPoolProvider := &fakebftdepl.FakeIpPoolProvider{}
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
		networksAssigner = NewSeededNetworksAssigner(networks, nameGenerator, ipPoolProvider, staticIpDecider, 5)
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
									IpPool:            expectedIpPool,
									AvailabilityZones: []string{"z1"},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpPool: expectedIpPool,
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
