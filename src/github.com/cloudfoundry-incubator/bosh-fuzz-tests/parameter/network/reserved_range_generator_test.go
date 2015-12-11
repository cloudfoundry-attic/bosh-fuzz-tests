package network_test

import (
	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReservedRangeGenerator", func() {
	var (
		generator ReservedRangeGenerator
	)

	BeforeEach(func() {
		decider := &fakebftdecider.FakeDecider{}
		decider.IsYesYes = true
		generator = NewReservedRangeGenerator("192.168.0", decider)
	})

	Describe("Generate", func() {
		It("generates ranges based on usedIps and reservedBorders", func() {
			usedIps := []int{15, 45, 75, 105}
			reservedBorders := []int{5, 60, 85, 130}
			reservedRanges := generator.Generate(usedIps, reservedBorders)
			Expect(reservedRanges).To(Equal([]string{
				"192.168.0.5-192.168.0.14",
				"192.168.0.16-192.168.0.44",
				"192.168.0.46-192.168.0.60",
				"192.168.0.85-192.168.0.104",
				"192.168.0.106-192.168.0.130",
			}))
		})

		It("generates correct ranges when usedIps are next to borders", func() {
			usedIps := []int{15, 45, 75, 76, 105}
			reservedBorders := []int{14, 46, 74, 77, 106}
			reservedRanges := generator.Generate(usedIps, reservedBorders)
			Expect(reservedRanges).To(Equal([]string{
				"192.168.0.14",
				"192.168.0.16-192.168.0.44",
				"192.168.0.46",
				"192.168.0.74",
				"192.168.0.77",
				"192.168.0.106",
			}))
		})
	})
})
