package deployment

import (
	"fmt"
)

type IpPool struct {
	IpRange string
	Gateway string
}

type IpPoolProvider interface {
	NewIpPool() IpPool
}

type ipPoolProvider struct {
	called int
}

func NewIpPoolProvider() IpPoolProvider {
	return &ipPoolProvider{}
}

func (p *ipPoolProvider) NewIpPool() IpPool {
	ipRange := fmt.Sprintf("192.168.%d.0/24", p.called)
	gateway := fmt.Sprintf("192.168.%d.1", p.called)
	p.called += 1

	return IpPool{
		IpRange: ipRange,
		Gateway: gateway,
	}
}
