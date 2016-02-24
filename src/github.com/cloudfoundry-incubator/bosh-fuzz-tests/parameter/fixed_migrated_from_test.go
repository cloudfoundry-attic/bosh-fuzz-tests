package parameter_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FixedMigratedFrom", func() {
	var (
		fixedMigratedFrom Parameter
	)

	BeforeEach(func() {
		fixedMigratedFrom = NewFixedMigratedFrom()
	})

	Context("when previous input does not have azs and current input has azs", func() {
		Context("when they have same job that is using same static IP", func() {
			It("specifies migrated_from on a job with az to which that static IP belongs", func() {
				input := bftinput.Input{
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
				previousInput := bftinput.Input{
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
				result := fixedMigratedFrom.Apply(input, previousInput)

				Expect(result.Jobs[0]).To(Equal(bftinput.Job{
					Name:              "foo-job",
					AvailabilityZones: []string{"z1", "z2"},
					MigratedFrom: []bftinput.MigratedFromConfig{
						{
							Name:             "foo-job",
							AvailabilityZone: "z2",
						},
					},
					Networks: []bftinput.JobNetworkConfig{
						{
							Name: "foo-network",
							StaticIps: []string{
								"192.168.2.232",
							},
						},
					},
				}))
			})
		})

		Context("when current job does not have any azs", func() {
			It("specifies migrated_from on a job with az to which that static IP belongs", func() {
				input := bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						Networks: []bftinput.NetworkConfig{
							{
								Name: "foo-network",
								Subnets: []bftinput.SubnetConfig{
									{
										IpPool: bftinput.NewIpPool("192.168.1", 1, []string{}),
									},
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
				result := fixedMigratedFrom.Apply(input, previousInput)

				Expect(result.Jobs[0]).To(Equal(bftinput.Job{
					Name: "foo-job",
					Networks: []bftinput.JobNetworkConfig{
						{
							Name: "foo-network",
							StaticIps: []string{
								"192.168.2.232",
							},
						},
					},
				}))
			})
		})
	})
})
