package network_test

import (
	"math/rand"

	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	fakebftnetwork "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network/fakes"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

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
		rand.Seed(64)

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
		ipPoolProvider.RegisterIpPool(ipPool)

		expectedIpPool = ipPool
		// reserving 2 ips since we have 2 instances
		expectedIpPool.NextStaticIp()
		expectedIpPool.NextStaticIp()

		decider = &fakebftdecider.FakeDecider{}
		decider.IsYesYes = true
		logger := boshlog.NewLogger(boshlog.LevelNone)
		networksAssigner = NewAssigner(networks, nameGenerator, ipPoolProvider, decider, logger)
	})

	It("assigns network of the given type to instance group and cloud config", func() {
		input := bftinput.Input{
			InstanceGroups: []bftinput.InstanceGroup{
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

		result := networksAssigner.Assign(input, bftinput.Input{})

		Expect(result).To(BeEquivalentTo(bftinput.Input{
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:              "foo",
					Instances:         2,
					AvailabilityZones: []string{"z1"},
					Networks: []bftinput.InstanceGroupNetworkConfig{
						{
							Name:          "foo-net",
							DefaultDNSnGW: true,
							StaticIps:     []string{"192.168.0.252", "192.168.0.219"},
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
				Compilation: bftinput.CompilationConfig{
					Network:          "foo-net",
					AvailabilityZone: "z1",
				},
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
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []bftinput.InstanceGroupNetworkConfig{
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

			result := networksAssigner.Assign(input, bftinput.Input{})

			Expect(result).To(BeEquivalentTo(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name:          "prev-net",
								DefaultDNSnGW: true,
								StaticIps:     []string{"192.168.0.252", "192.168.0.219"},
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
							},
						},
						{
							Name: "baz-net",
							Type: "vip",
						},
					},
					Compilation: bftinput.CompilationConfig{
						Network:          "prev-net",
						AvailabilityZone: "z1",
					},
				},
			},
			))
		})
	})

	It("does not reuse IP if IP does not belong to network subnet", func() {
		input := bftinput.Input{
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:      "foo",
					Instances: 2,
					Networks: []bftinput.InstanceGroupNetworkConfig{
						{
							Name: "default",
						},
					},
				},
			},
			CloudConfig: bftinput.CloudConfig{
				Networks: []bftinput.NetworkConfig{
					{
						Name: "default",
						Type: "manual",
						Subnets: []bftinput.SubnetConfig{
							{
								IpPool: bftinput.NewIpPool("192.168.0", 1, []string{}),
							},
						},
					},
				},
			},
		}

		previousInput := bftinput.Input{
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:      "foo",
					Instances: 2,
					Networks: []bftinput.InstanceGroupNetworkConfig{
						{
							Name:      "default",
							StaticIps: []string{"192.168.4.209", "192.168.4.254"},
						},
					},
				},
			},
			CloudConfig: bftinput.CloudConfig{
				Networks: []bftinput.NetworkConfig{
					{
						Name: "default",
						Type: "manual",
						Subnets: []bftinput.SubnetConfig{
							{
								IpPool: bftinput.NewIpPool("192.168.4", 1, []string{}),
							},
						},
					},
				},
			},
		}

		result := networksAssigner.Assign(input, previousInput)
		Expect(len(result.InstanceGroups[0].Networks[0].StaticIps)).To(Equal(2))
		Expect(result.InstanceGroups[0].Networks[0].StaticIps).ToNot(ContainElement("192.168.4.209"))
		Expect(result.InstanceGroups[0].Networks[0].StaticIps).ToNot(ContainElement("192.168.4.254"))
	})

	Context("when previous input has static IPs", func() {
		BeforeEach(func() {
			decider.IsYesYes = true

			expectedIpPool.ReserveStaticIp("192.168.0.219")
			expectedIpPool.ReserveStaticIp("192.168.0.245")
			expectedIpPool.ReserveStaticIp("192.168.0.252")
			expectedIpPool.ReserveStaticIp("192.168.0.236")
		})

		It("does not reuse those static IPs", func() {
			// our fuzzing returns ips in order
			// "192.168.0.252", "192.168.0.219", "192.168.0.234", "192.168.0.245"
			// we put "192.168.0.252" on second instanceGroup to make sure it is not going to be used by first instanceGroup
			// we put "192.168.0.245" on first instanceGroup to make sure it is not going to be used by second instanceGroup
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:      "foo",
						Instances: 2,
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name:      "prev-net",
								StaticIps: []string{"192.168.0.219", "192.168.0.245"},
							},
						},
					},
					{
						Name:      "bar",
						Instances: 2,
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name:      "prev-net",
								StaticIps: []string{"192.168.0.252", "192.168.0.236"},
							},
						},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					Networks: []bftinput.NetworkConfig{
						{
							Name: "prev-net",
							Type: "manual",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: expectedIpPool,
								},
							},
						},
					},
				},
			}

			result := networksAssigner.Assign(input, input)

			Expect(result).To(BeEquivalentTo(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:      "foo",
						Instances: 2,
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name:          "prev-net",
								DefaultDNSnGW: true,
								StaticIps:     []string{"192.168.0.234", "192.168.0.223"},
							},
						},
					},
					{
						Name:      "bar",
						Instances: 2,
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name:          "prev-net",
								DefaultDNSnGW: true,
								StaticIps:     []string{"192.168.0.214", "192.168.0.228"},
							},
						},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					Networks: []bftinput.NetworkConfig{
						{
							Name: "prev-net",
							Type: "manual",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: expectedIpPool,
								},
							},
						},
						{
							Name: "foo-net",
							Type: "vip",
						},
					},
					Compilation: bftinput.CompilationConfig{
						Network: "prev-net",
					},
				},
			},
			))
		})
	})
})
