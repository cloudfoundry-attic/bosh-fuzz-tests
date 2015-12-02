package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeIpPoolProvider struct {
	IpPools []*bftinput.IpPool
}

func (f *FakeIpPoolProvider) NewIpPool(numOfNeededIPs int) *bftinput.IpPool {
	var ipPool *bftinput.IpPool
	ipPool, f.IpPools = f.IpPools[0], f.IpPools[1:]
	return ipPool
}

func (f *FakeIpPoolProvider) Reset() {}

func (f *FakeIpPoolProvider) RegisterIpPool(ipPool *bftinput.IpPool) {
	f.IpPools = append(f.IpPools, ipPool)
}
