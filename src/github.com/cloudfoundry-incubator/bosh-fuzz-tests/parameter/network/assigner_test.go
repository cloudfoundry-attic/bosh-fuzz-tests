package network_test

import (
	"math/rand"

	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	fakebftnetwork "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworksAssigner", func() {
	var (
		networksAssigner Assigner
		networks         [][]string
		expectedIpPool   *bftinput.IpPool
		decider          *fakebftdecider.FakeDecider
	)

	BeforeEach(func() {
		rand.Seed(32)

		networks = [][]string{[]string{"manual", "vip"}}
		nameGenerator := &fakebftnamegen.FakeNameGenerator{}
		nameGenerator.Names = []string{"foo-net", "bar-net", "baz-net", "qux-net"}
		ipPoolProvider := &fakebftnetwork.FakeIpPoolProvider{}
		vipPool := &bftinput.IpPool{}
		ipPoolProvider.RegisterIpPool(vipPool)

		ipPool := bftinput.NewIpPool(
			"192.168.0",
			1,
			[]string{
				"192.168.0.15-192.168.0.58",
				"192.168.0.157-192.168.0.203",
			},
		)
		ipPoolProvider.RegisterIpPool(ipPool)
		ipPoolProvider.RegisterIpPool(ipPool)

		ipPoolProvider.RegisterIpPool(ipPool)

		expectedIpPool = bftinput.NewIpPool(
			"192.168.0",
			1,
			[]string{
				"192.168.0.15-192.168.0.58",
				"192.168.0.157-192.168.0.203",
			},
		)
		// reserving 2 ips since we have 2 instances
		expectedIpPool.NextStaticIp()
		expectedIpPool.NextStaticIp()

		decider = &fakebftdecider.FakeDecider{}
		decider.IsYesYes = true
		networksAssigner = NewAssigner(networks, nameGenerator, ipPoolProvider, decider)
	})

	It("assigns network of the given type to job and cloud config", func() {
		input := bftinput.Input{
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
		}

		result := networksAssigner.Assign(input)

		Expect(result).To(BeEquivalentTo(bftinput.Input{
			Jobs: []bftinput.Job{
				{
					Name:              "foo",
					Instances:         2,
					AvailabilityZones: []string{"z1"},
					Networks: []bftinput.JobNetworkConfig{
						{
							Name:          "foo-net",
							DefaultDNSnGW: true,
							StaticIps:     []string{"192.168.0.200", "192.168.0.201"},
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
		))
	})

	Context("when it is decided to reuse same network name", func() {
		BeforeEach(func() {
			decider.IsYesYes = true
		})

		It("reuses network name from previous input", func() {
			input := bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []bftinput.JobNetworkConfig{
							{Name: "prev-net"},
						},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
					Networks: []bftinput.NetworkConfig{
						{
							Name: "prev-net",
							Type: "dynamic",
						},
					},
				},
			}

			result := networksAssigner.Assign(input)

			Expect(result).To(BeEquivalentTo(bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []bftinput.JobNetworkConfig{
							{
								Name:          "prev-net",
								DefaultDNSnGW: true,
								StaticIps:     []string{"192.168.0.200", "192.168.0.201"},
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
							Name: "prev-net",
							Type: "manual",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool:            expectedIpPool,
									AvailabilityZones: []string{"z1"},
								},
							},
						},
						{
							Name: "foo-net",
							Type: "vip",
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: expectedIpPool,
								},
								{
									IpPool: expectedIpPool,
								},
							},
						},
						{
							Name: "baz-net",
							Type: "vip",
						},
					},
					CompilationNetwork: "bar-net",
				},
			},
			))
		})
	})
})
