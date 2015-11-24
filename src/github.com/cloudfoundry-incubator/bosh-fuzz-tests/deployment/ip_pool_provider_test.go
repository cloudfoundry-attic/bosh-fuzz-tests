package deployment_test

import (
	"math/rand"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

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

		It("generates list of reserved IPs", func() {
			ipPool := ipPoolProvider.NewIpPool(5)
			Expect(ipPool.Reserved).To(Equal([]string{
				"192.168.0.41-192.168.0.49",
				"192.168.0.51-192.168.0.104",
				"192.168.0.173",
				"192.168.0.209",
				"192.168.0.242",
			}))
		})

		It("generates list of available IPs", func() {
			ipPool := ipPoolProvider.NewIpPool(3)
			ip, err := ipPool.NextStaticIp()
			Expect(err).ToNot(HaveOccurred())
			Expect(ip).To(Equal("192.168.0.207"))

			ip, err = ipPool.NextStaticIp()
			Expect(err).ToNot(HaveOccurred())
			Expect(ip).To(Equal("192.168.0.206"))

			ip, err = ipPool.NextStaticIp()
			Expect(err).ToNot(HaveOccurred())
			Expect(ip).To(Equal("192.168.0.4"))

			ip, err = ipPool.NextStaticIp()
			Expect(err).To(HaveOccurred())
		})
	})
})
