package analyzer_test

import (
	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analyzer", func() {
	var (
		analyzer bftanalyzer.Analyzer
	)

	BeforeEach(func() {
		analyzer = bftanalyzer.NewAnalyzer(nil)
	})

	Context("when previous input has azs and current input does not have azs", func() {
		Context("when they have the same job that is using the same static IP", func() {
			It("specifies migrated_from on a job without an azs", func() {
				input := bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						AvailabilityZones: nil,
						Networks: []bftinput.NetworkConfig{
							{
								Name: "foo-network",
								Subnets: []bftinput.SubnetConfig{
									{
										IpPool: bftinput.NewIpPool("192.168.2", 1, []string{}),
									},
								},
							},
						},
					},
					Jobs: []bftinput.Job{
						{
							Name: "foo-job",
							Networks: []bftinput.JobNetworkConfig{
								{
									Name: "foo-network",
									StaticIps: []string{
										"192.168.2.232",
									},
								},
							},
						},
					},
				}
				previousInput := bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						AvailabilityZones: []bftinput.AvailabilityZone{
							{
								Name: "z1",
							},
							{
								Name: "z2",
							},
						},
						Networks: []bftinput.NetworkConfig{
							{
								Name: "foo-network",
								Subnets: []bftinput.SubnetConfig{
									{
										AvailabilityZones: []string{"z1"},
										IpPool:            bftinput.NewIpPool("192.168.1", 1, []string{}),
									},
									{
										AvailabilityZones: []string{"z2"},
										IpPool:            bftinput.NewIpPool("192.168.2", 1, []string{}),
									},
								},
							},
						},
					},
					Jobs: []bftinput.Job{
						{
							Name:              "foo-job",
							AvailabilityZones: []string{"z1", "z2"},
							Networks: []bftinput.JobNetworkConfig{
								{
									Name: "foo-network",
									StaticIps: []string{
										"192.168.2.232",
									},
								},
							},
						},
					},
				}

				result := analyzer.Analyze([]bftinput.Input{previousInput, input})

				Expect(result[0].DeploymentWillFail).To(BeFalse())
				Expect(result[1].DeploymentWillFail).To(BeTrue())
			})
		})
	})

	It("It does not move an existing instance's static IP to another AZ", func() {
		previousInput := bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []bftinput.AvailabilityZone{
					{
						Name: "z1",
					},
				},
				Networks: []bftinput.NetworkConfig{
					{
						Name: "foo-network",
						Subnets: []bftinput.SubnetConfig{
							{
								AvailabilityZones: []string{"z1"},
								IpPool:            bftinput.NewIpPool("192.168.2", 1, []string{}),
							},
						},
					},
				},
			},
			Jobs: []bftinput.Job{
				{
					Name:              "foo-job",
					AvailabilityZones: []string{"z1"},
					Networks: []bftinput.JobNetworkConfig{
						{
							Name: "foo-network",
							StaticIps: []string{
								"192.168.2.232",
							},
						},
					},
				},
			},
		}

		input := bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []bftinput.AvailabilityZone{
					{
						Name: "z2",
					},
				},
				Networks: []bftinput.NetworkConfig{
					{
						Name: "foo-network",
						Subnets: []bftinput.SubnetConfig{
							{
								AvailabilityZones: []string{"z2"},
								IpPool:            bftinput.NewIpPool("192.168.2", 1, []string{}),
							},
						},
					},
				},
			},
			Jobs: []bftinput.Job{
				{
					Name:              "foo-job",
					AvailabilityZones: []string{"z2"},
					Networks: []bftinput.JobNetworkConfig{
						{
							Name: "foo-network",
							StaticIps: []string{
								"192.168.2.232",
							},
						},
					},
				},
			},
		}

		result := analyzer.Analyze([]bftinput.Input{previousInput, input})

		Expect(result[0].DeploymentWillFail).To(BeFalse())
		Expect(result[1].DeploymentWillFail).To(BeTrue())
	})
})
