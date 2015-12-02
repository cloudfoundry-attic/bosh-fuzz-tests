package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type stemcell struct {
	definition string
}

func NewStemcell(definition string) Parameter {
	return &stemcell{
		definition: definition,
	}
}

func (s *stemcell) Apply(input *bftinput.Input) *bftinput.Input {
	var stemcellConfig bftinput.StemcellConfig

	if s.definition == "os" {
		stemcellConfig = bftinput.StemcellConfig{
			OS:      "toronto-os",
			Version: "1",
		}
	} else {
		stemcellConfig = bftinput.StemcellConfig{
			Name:    "ubuntu-stemcell",
			Version: "1",
		}
	}

	if len(input.CloudConfig.VmTypes) > 0 {
		stemcellConfig.Alias = "default"
		input.Stemcells = []bftinput.StemcellConfig{stemcellConfig}
	} else {
		for r, _ := range input.CloudConfig.ResourcePools {
			input.CloudConfig.ResourcePools[r].Stemcell = stemcellConfig
		}
	}

	return input
}
