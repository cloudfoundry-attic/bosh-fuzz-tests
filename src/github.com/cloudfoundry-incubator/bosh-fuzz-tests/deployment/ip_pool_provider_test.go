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
			ipPool := ipPoolProvider.NewIpPool()
			Expect(ipPool.IpRange).To(Equal("192.168.0.0/24"))

			ipPool = ipPoolProvider.NewIpPool()
			Expect(ipPool.IpRange).To(Equal("192.168.1.0/24"))

			ipPool = ipPoolProvider.NewIpPool()
			Expect(ipPool.IpRange).To(Equal("192.168.2.0/24"))
		})
		It("alternates gateway's last IP between .1 and .254", func() {
			ipPool := ipPoolProvider.NewIpPool()
			Expect(ipPool.IpRange).To(Equal("192.168.0.0/24"))
		})

		It("generates list of reserved IPs", func() {
			ipPool := ipPoolProvider.NewIpPool()
			Expect(ipPool.Reserved).To(Equal([]string{
				"192.168.0.3",
				"192.168.0.54-192.168.0.221",
				"192.168.0.223-192.168.0.242",
			}))
		})
	})
})
