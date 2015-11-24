package deployment

import "fmt"

type ReservedRangeGenerator interface {
	Generate(usedIps []int, reservedBorders []int) []string
}

type reservedRangeGenerator struct {
	prefix  string
	decider Decider
}

func NewReservedRangeGenerator(prefix string, decider Decider) ReservedRangeGenerator {
	return &reservedRangeGenerator{
		prefix:  prefix,
		decider: decider,
	}
}

func (r *reservedRangeGenerator) Generate(usedIps []int, reservedBorders []int) []string {
	reservedRanges := []string{}
	var currentBorder, nextBorder int

	for len(reservedBorders) > 0 {
		currentBorder, reservedBorders = reservedBorders[0], reservedBorders[1:]
		if r.decider.IsYes() && len(reservedBorders) > 0 {
			nextBorder, reservedBorders = reservedBorders[0], reservedBorders[1:]
			firstIpInRange := currentBorder
			for _, usedIp := range usedIps {
				if firstIpInRange == usedIp {
					firstIpInRange = usedIp + 1
					continue
				}
				if usedIp > firstIpInRange && usedIp < nextBorder {
					if firstIpInRange == usedIp-1 {
						reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d", r.prefix, firstIpInRange))
					} else {
						reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d-%s.%d", r.prefix, firstIpInRange, r.prefix, usedIp-1))
					}
					firstIpInRange = usedIp + 1
				}
			}

			if firstIpInRange < nextBorder {
				reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d-%s.%d", r.prefix, firstIpInRange, r.prefix, nextBorder))
			} else if firstIpInRange == nextBorder {
				reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d", r.prefix, nextBorder))
			}
		} else {
			reservedRanges = append(reservedRanges, fmt.Sprintf("%s.%d", r.prefix, currentBorder))
		}
	}

	return reservedRanges
}
