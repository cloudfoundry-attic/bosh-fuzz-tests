package deployment_test

import (
	"math/rand"

	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	fakebftdepl "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment/fakes"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworksAssigner", func() {
	var (
		networksAssigner NetworksAssigner
		networks         [][]string
		expectedIpPool   *bftinput.IpPool
	)

	BeforeEach(func() {
		rand.Seed(32)

		networks = [][]string{[]string{"manual", "vip"}}
		nameGenerator := &fakebftnamegen.FakeNameGenerator{}
		nameGenerator.Names = []string{"foo-net", "bar-net", "baz-net", "qux-net"}
		ipPoolProvider := &fakebftdepl.FakeIpPoolProvider{}
		vipPool := &bftinput.IpPool{
			AvailableIps: []string{
				"10.10.0.6",
				"10.10.0.32",
			},
		}
		ipPoolProvider.RegisterIpPool(vipPool)

		ipPool := &bftinput.IpPool{
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

		expectedIpPool = &bftinput.IpPool{
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
		staticIpDecider := &fakebftdecider.FakeDecider{}
		staticIpDecider.IsYesYes = true
		networksAssigner = NewNetworksAssigner(networks, nameGenerator, ipPoolProvider, staticIpDecider)
	})

	It("assigns network of the given type to job and cloud config", func() {
		inputs := []bftinput.Input{
			{
				Jobs: []bftinput.Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
				},
			},
		}

		networksAssigner.Assign(inputs)

		Expect(inputs).To(BeEquivalentTo([]bftinput.Input{
			{
				Jobs: []bftinput.Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []bftinput.JobNetworkConfig{
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
								StaticIps:     []string{"192.168.0.222", "192.168.0.110"},
							},
						},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
					Networks: []bftinput.NetworkConfig{
						{
							Name: "foo-net",
							Type: "manual",
							Subnets: []bftinput.SubnetConfig{
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
							Subnets: []bftinput.SubnetConfig{
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
