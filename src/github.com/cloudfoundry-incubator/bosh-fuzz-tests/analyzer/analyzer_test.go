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
		Context("when they have the same instance group that is using the same static IP", func() {
			It("specifies migrated_from on an instance group without an azs", func() {
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
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name: "foo-instance-group",
							Networks: []bftinput.InstanceGroupNetworkConfig{
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
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:              "foo-instance-group",
							AvailabilityZones: []string{"z1", "z2"},
							Networks: []bftinput.InstanceGroupNetworkConfig{
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

	Context("when variable is a certificate", func() {
		Context("and certificate is not a CA", func() {
			Context("and it does not reference a signing CA", func() {
				It("should expect deployment to fail", func() {
					input := bftinput.Input{
						Variables: []bftinput.Variable{
							{
								Name:    "bad_cert",
								Type:    "certificate",
								Options: map[string]interface{}{"is_ca": false},
							},
						},
					}
					result := analyzer.Analyze([]bftinput.Input{input})
					Expect(result[0].DeploymentWillFail).To(BeTrue())
				})
			})
			Context("and it references a non-existent variable as the signing CA", func() {
				It("should expect deployment to fail", func() {
					input := bftinput.Input{
						Variables: []bftinput.Variable{
							{
								Name:    "bad_cert",
								Type:    "certificate",
								Options: map[string]interface{}{"is_ca": false, "ca": "nonexistent"},
							},
						},
					}
					result := analyzer.Analyze([]bftinput.Input{input})
					Expect(result[0].DeploymentWillFail).To(BeTrue())
				})
			})
			Context("and it references a non-CA variable as the signing CA", func() {
				It("should expect deployment to fail", func() {
					input := bftinput.Input{
						Variables: []bftinput.Variable{
							{
								Name:    "signed_cert",
								Type:    "certificate",
								Options: map[string]interface{}{"is_ca": false, "ca": "bad_signing_cert"},
							},
							{
								Name:    "bad_signing_cert",
								Type:    "certificate",
								Options: map[string]interface{}{"is_ca": false, "ca": "root_ca"},
							},
							{
								Name:    "root_ca",
								Type:    "certificate",
								Options: map[string]interface{}{"is_ca": true},
							},
						},
					}
					result := analyzer.Analyze([]bftinput.Input{input})
					Expect(result[0].DeploymentWillFail).To(BeTrue())
				})
			})
		})

		Context("and certificate is a CA", func() {
			Context("and it references another certificate", func() {
				Context("and the referenced certificate is a CA", func() {
					It("should expect deployment to not fail", func() {
						input := bftinput.Input{
							Variables: []bftinput.Variable{
								{
									Name:    "signed_cert",
									Type:    "certificate",
									Options: map[string]interface{}{"is_ca": true, "ca": "root_ca"},
								},
								{
									Name:    "root_ca",
									Type:    "certificate",
									Options: map[string]interface{}{"is_ca": true},
								},
							},
						}
						result := analyzer.Analyze([]bftinput.Input{input})
						Expect(result[0].DeploymentWillFail).To(BeFalse())
					})
				})
				Context("and the referenced certificate is missing", func() {
					It("should expect deployment to fail", func() {
						input := bftinput.Input{
							Variables: []bftinput.Variable{
								{
									Name:    "signed_cert",
									Type:    "certificate",
									Options: map[string]interface{}{"is_ca": true, "ca": "root_ca"},
								},
							},
						}
						result := analyzer.Analyze([]bftinput.Input{input})
						Expect(result[0].DeploymentWillFail).To(BeTrue())
					})
				})
				Context("and the referenced certificate is not a CA", func() {
					It("should expect deployment to fail", func() {
						input := bftinput.Input{
							Variables: []bftinput.Variable{
								{
									Name:    "signed_cert",
									Type:    "certificate",
									Options: map[string]interface{}{"is_ca": true, "ca": "bad_signing_cert"},
								},
								{
									Name:    "bad_signing_cert",
									Type:    "certificate",
									Options: map[string]interface{}{"is_ca": false, "ca": "root_ca"},
								},
								{
									Name:    "root_ca",
									Type:    "certificate",
									Options: map[string]interface{}{"is_ca": true},
								},
							},
						}
						result := analyzer.Analyze([]bftinput.Input{input})
						Expect(result[0].DeploymentWillFail).To(BeTrue())
					})
				})
			})
		})
	})

	Context("when some of the inputs are dry-run", func() {
		It("only considers non dry-run inputs when building expectations", func() {
			normalInput1 := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "ig1",
					},
				},
			}
			normalInput2 := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "ig1",
					},
				},
			}
			dryRunInput := bftinput.Input{
				IsDryRun: true,
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "other-ig",
					},
				},
			}

			result := analyzer.Analyze([]bftinput.Input{normalInput1, dryRunInput, normalInput2})
			Expect(result[1].Expectations).To(BeEmpty())
			Expect(result[2].Expectations).To(HaveLen(2))
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
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:              "foo-instance-group",
					AvailabilityZones: []string{"z1"},
					Networks: []bftinput.InstanceGroupNetworkConfig{
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
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:              "foo-instance-group",
					AvailabilityZones: []string{"z2"},
					Networks: []bftinput.InstanceGroupNetworkConfig{
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

	It("does not move an existing instance's static IP to another instance group", func() {
		previousInput := bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
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
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name: "foo-instance-group",
					Networks: []bftinput.InstanceGroupNetworkConfig{
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
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name: "foo-instance-group-renamed",
					Networks: []bftinput.InstanceGroupNetworkConfig{
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
