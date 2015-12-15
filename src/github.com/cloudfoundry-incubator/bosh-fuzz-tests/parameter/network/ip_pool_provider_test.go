package network_test

import (
	"math/rand"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IpPoolProvider", func() {
	var (
		ipPoolProvider IpPoolProvider
	)

	BeforeEach(func() {
		rand.Seed(65)
		ipPoolProvider = NewIpPoolProvider()
	})

	Describe("NewIpPool", func() {
		It("generates new pool on subsequent calls", func() {
			ipPool := ipPoolProvider.NewIpPool(1)
			Expect(ipPool.IpRange).To(Equal("192.168.0.0/24"))

			ipPool = ipPoolProvider.NewIpPool(1)
			Expect(ipPool.IpRange).To(Equal("192.168.1.0/24"))

			ipPool = ipPoolProvider.NewIpPool(1)
			Expect(ipPool.IpRange).To(Equal("192.168.2.0/24"))
		})

		It("alternates gateway's last IP between .1 and .254", func() {
			ipPool := ipPoolProvider.NewIpPool(1)
			Expect(ipPool.IpRange).To(Equal("192.168.0.0/24"))
		})

		It("generates static IP list", func() {
			ipPool := ipPoolProvider.NewIpPool(5)
			Expect(ipPool.Static).To(Equal([]string{
				"192.168.0.200-192.168.0.253",
			}))
		})

		It("generates list of reserved IPs", func() {
			ipPool := ipPoolProvider.NewIpPool(5)
			Expect(ipPool.Reserved).To(Equal([]string{
				"192.168.0.41",
				"192.168.0.68-192.168.0.104",
				"192.168.0.168-192.168.0.173",
			}))
		})

		It("generates list of static IPs", func() {
			ipPool := ipPoolProvider.NewIpPool(1)
			ip, err := ipPool.NextStaticIp()
			Expect(err).ToNot(HaveOccurred())
			Expect(ip).To(Equal("192.168.0.215"))

			ip, err = ipPool.NextStaticIp()
			Expect(err).ToNot(HaveOccurred())
			Expect(ip).To(Equal("192.168.0.250"))

			ip, err = ipPool.NextStaticIp()
			Expect(err).ToNot(HaveOccurred())
			Expect(ip).To(Equal("192.168.0.208"))
		})
	})
})
