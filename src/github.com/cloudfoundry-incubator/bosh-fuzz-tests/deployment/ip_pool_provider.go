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
	called             int
	gatewayFourthOctet int
}

func NewIpPoolProvider() IpPoolProvider {
	return &ipPoolProvider{}
}

func (p *ipPoolProvider) NewIpPool() IpPool {
	if p.gatewayFourthOctet == 1 {
		p.gatewayFourthOctet = 254
	} else {
		p.gatewayFourthOctet = 1
	}

	ipRange := fmt.Sprintf("192.168.%d.0/24", p.called)
	gateway := fmt.Sprintf("192.168.%d.%d", p.called, p.gatewayFourthOctet)
	p.called += 1

	return IpPool{
		IpRange: ipRange,
		Gateway: gateway,
	}
}
