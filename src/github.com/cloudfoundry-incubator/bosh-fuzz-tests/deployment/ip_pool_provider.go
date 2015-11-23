package deployment

import (
	"fmt"
	"math/rand"
	"sort"
)

type IpPool struct {
	IpRange  string
	Gateway  string
	Reserved []string
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

	prefix := fmt.Sprintf("192.168.%d", p.called)
	p.called += 1

	ipRange := fmt.Sprintf("%s.0/24", prefix)
	gateway := fmt.Sprintf("%s.%d", prefix, p.gatewayFourthOctet)

	numberOfReservedBorders := rand.Intn(6) // up to 6 borders of reserved ranges
	reservedBorders := []int{}
	for i := 0; i < numberOfReservedBorders; i++ {
		reservedBorders = append(reservedBorders, rand.Intn(253)+p.gatewayFourthOctet%254+1) // skip 0 and gateway(1, 254)
	}

	sort.Ints(reservedBorders)

	reservedRanges := []string{}
	var currentBorder, nextBorder int
	for len(reservedBorders) > 0 {
		currentBorder, reservedBorders = reservedBorders[0], reservedBorders[1:]
		if rand.Intn(2) == 1 && len(reservedBorders) > 0 {
			nextBorder, reservedBorders = reservedBorders[0], reservedBorders[1:]
			reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d-%s.%d", prefix, currentBorder, prefix, nextBorder))
		} else {
			reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d", prefix, currentBorder))
		}
	}

	return IpPool{
		IpRange:  ipRange,
		Gateway:  gateway,
		Reserved: reservedRanges,
	}
}
