package analyzer_test

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StemcellComparator", func() {
	var (
		stemcellComparator Comparator
		previousInput      bftinput.Input
		currentInput       bftinput.Input
	)

	BeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelNone)
		stemcellComparator = NewStemcellComparator(logger)
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
			expectedDebugLogExpectation := bftexpectation.NewExistingInstanceDebugLog("stemcell_changed?")
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
			expectedDebugLogExpectation := bftexpectation.NewExistingInstanceDebugLog("stemcell_changed?")
			Expect(expectations).To(ContainElement(expectedDebugLogExpectation))
		})
	})
})
