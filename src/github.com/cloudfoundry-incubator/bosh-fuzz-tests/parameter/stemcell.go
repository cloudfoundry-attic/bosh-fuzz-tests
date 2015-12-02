package parameter

import (
	"fmt"
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type stemcell struct {
	definition       string
	stemcellVersions []string
}

func NewStemcell(definition string, stemcellVersions []string) Parameter {
	return &stemcell{
		definition:       definition,
		stemcellVersions: stemcellVersions,
	}
}

func (s *stemcell) Apply(input *bftinput.Input) *bftinput.Input {
	var stemcellConfig bftinput.StemcellConfig

	if s.definition == "os" {
		stemcellConfig = bftinput.StemcellConfig{
			OS: "toronto-os",
		}
	} else {
		stemcellConfig = bftinput.StemcellConfig{
			Name: "ubuntu-stemcell",
		}
	}

	if len(input.CloudConfig.VmTypes) > 0 {
		for _, vmType := range input.CloudConfig.VmTypes {
			stemcellConfig.Version = s.stemcellVersions[rand.Intn(len(s.stemcellVersions))]
			stemcellConfig.Alias = fmt.Sprintf("stemcell-%s", stemcellConfig.Version)
			input.Stemcells = append(input.Stemcells, stemcellConfig)
			for j := range input.Jobs {
				if input.Jobs[j].VmType == vmType.Name {
					input.Jobs[j].Stemcell = stemcellConfig.Alias
				}
			}
		}
	} else {
		for r, _ := range input.CloudConfig.ResourcePools {
			stemcellConfig.Version = s.stemcellVersions[rand.Intn(len(s.stemcellVersions))]
			input.CloudConfig.ResourcePools[r].Stemcell = stemcellConfig
		}
	}

	return input
}
