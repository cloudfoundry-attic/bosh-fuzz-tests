package analyzer_test

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StemcellComparator", func() {
	var (
		stemcellComparator Comparator
		previousInput      bftinput.Input
		currentInput       bftinput.Input
		cliRunner          bltclirunner.Runner
	)

	BeforeEach(func() {
		fs := fakesys.NewFakeFileSystem()
		cmdRunner := fakesys.NewFakeCmdRunner()
		boshCmd := boshsys.Command{Name: "bosh"}

		cliRunner = bltclirunner.NewRunner(boshCmd, cmdRunner, fs)
		cliRunner.Configure()

		expectationFactory := bftexpectation.NewFactory(cliRunner)
		stemcellComparator = NewStemcellComparator(expectationFactory)
	})

	Context("when there are same jobs that have different stemcell versions using vm types", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Stemcells: []bftinput.StemcellConfig{
					{
						Alias:   "fake-stemcell",
						Version: "1",
					},
				},
				Jobs: []bftinput.Job{
					{
						Name:     "foo-job",
						Stemcell: "fake-stemcell",
					},
				},
			}

			currentInput = bftinput.Input{
				Stemcells: []bftinput.StemcellConfig{
					{
						Alias:   "fake-stemcell",
						Version: "2",
					},
				},
				Jobs: []bftinput.Job{
					{
						Name:     "foo-job",
						Stemcell: "fake-stemcell",
					},
				},
			}
		})

		It("returns debug log expectation", func() {
			expectations := stemcellComparator.Compare(previousInput, currentInput)
			expectedDebugLogExpectation := bftexpectation.NewDebugLog("stemcell_changed?", cliRunner)
			Expect(expectations).To(ContainElement(expectedDebugLogExpectation))
		})
	})

	Context("when there are same jobs that have different stemcell versions using resource pools", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name: "fake-resource-pool",
							Stemcell: bftinput.StemcellConfig{
								Name:    "fake-stemcell",
								Version: "1",
							},
						},
					},
				},
				Jobs: []bftinput.Job{
					{
						Name:         "foo-job",
						ResourcePool: "fake-resource-pool",
					},
				},
			}

			currentInput = bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name: "fake-resource-pool",
							Stemcell: bftinput.StemcellConfig{
								Name:    "fake-stemcell",
								Version: "2",
							},
						},
					},
				},
				Jobs: []bftinput.Job{
					{
						Name:         "foo-job",
						ResourcePool: "fake-resource-pool",
					},
				},
			}
		})

		It("returns debug log expectation", func() {
			expectations := stemcellComparator.Compare(previousInput, currentInput)
			expectedDebugLogExpectation := bftexpectation.NewDebugLog("stemcell_changed?", cliRunner)
			Expect(expectations).To(ContainElement(expectedDebugLogExpectation))
		})
	})
})
