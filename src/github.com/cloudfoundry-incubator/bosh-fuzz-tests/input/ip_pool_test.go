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

	Describe("Contains", func() {
		BeforeEach(func() {
			ipPool = NewIpPool(
				"10.10.2",
				1,
				[]string{},
			)
		})
		It("can tell when an IP address is within a subnet", func() {
			Expect(ipPool.Contains("10.10.2.10")).To(BeTrue())
			Expect(ipPool.Contains("10.10.2.1")).To(BeTrue())
			Expect(ipPool.Contains("10.10.2.254")).To(BeTrue())
		})
		It("can tell when an IP address is NOT within a subnet", func() {
			Expect(ipPool.Contains("10.10.20.20")).To(BeFalse())
			Expect(ipPool.Contains("10.10.255.10")).To(BeFalse())
			Expect(ipPool.Contains("10.10.30.20")).To(BeFalse())
			Expect(ipPool.Contains("192.168.2.254")).To(BeFalse())
			Expect(ipPool.Contains("224.10.20.20")).To(BeFalse())
		})
	})
})
