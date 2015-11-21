package deployment_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IpPoolProvider", func() {
	var (
		ipPoolProvider IpPoolProvider
	)

	BeforeEach(func() {
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
	})
})
