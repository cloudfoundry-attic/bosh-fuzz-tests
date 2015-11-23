package deployment

type PlacedAZs struct {
	azs map[string]bool
}

func NewPlacedAZs() *PlacedAZs {
	return &PlacedAZs{
		azs: map[string]bool{},
	}
}

func (a *PlacedAZs) Place(azs []string) {
	for _, az := range azs {
		a.azs[az] = true
	}
}

func (a *PlacedAZs) AllPlaced(azs []string) bool {
	for _, az := range azs {
		if !a.azs[az] {
			return false
		}
	}

	return true
}
