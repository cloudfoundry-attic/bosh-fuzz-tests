package input_test

import (
	"math/rand"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IpPool", func() {
	var (
		ipPool *IpPool
	)

	BeforeEach(func() {
		rand.Seed(64)

		ipPool = NewIpPool(
			"10.10.0",
			1,
			[]string{},
		)
	})

	Describe("NextStaticIp", func() {
		It("returns the next static IP in static range", func() {
			staticIp, err := ipPool.NextStaticIp()
			Expect(err).NotTo(HaveOccurred())
			Expect(staticIp).To(Equal("10.10.0.241"))

			staticIp, err = ipPool.NextStaticIp()
			Expect(err).NotTo(HaveOccurred())
			Expect(staticIp).To(Equal("10.10.0.249"))
		})
	})

	Describe("ReserveStaticIp", func() {
		It("returns the next static IP in static range", func() {
			ipPool.ReserveStaticIp("10.10.0.241")

			staticIp, err := ipPool.NextStaticIp()
			Expect(err).NotTo(HaveOccurred())
			Expect(staticIp).To(Equal("10.10.0.249"))
		})
	})
})
