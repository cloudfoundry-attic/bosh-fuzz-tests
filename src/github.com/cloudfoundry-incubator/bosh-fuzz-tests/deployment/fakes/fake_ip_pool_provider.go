package fakes

import (
	bftdeployment "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
)

type FakeIpPoolProvider struct {
	IpPools []*bftdeployment.IpPool
}

func (f *FakeIpPoolProvider) NewIpPool(numOfNeededIPs int) *bftdeployment.IpPool {
	var ipPool *bftdeployment.IpPool
	ipPool, f.IpPools = f.IpPools[0], f.IpPools[1:]
	return ipPool
}

func (f *FakeIpPoolProvider) Reset() {}

func (f *FakeIpPoolProvider) RegisterIpPool(ipPool *bftdeployment.IpPool) {
	f.IpPools = append(f.IpPools, ipPool)
}
