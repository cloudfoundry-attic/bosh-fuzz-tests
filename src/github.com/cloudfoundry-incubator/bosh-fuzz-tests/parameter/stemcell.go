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

func (s *stemcell) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.Stemcells = nil

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

	usedStemcells := map[string]bool{}

	if len(input.CloudConfig.VmTypes) > 0 {
		for _, vmType := range input.CloudConfig.VmTypes {
			stemcellConfig.Version = s.stemcellVersions[rand.Intn(len(s.stemcellVersions))]
			stemcellConfig.Alias = fmt.Sprintf("stemcell-%s", stemcellConfig.Version)

			if usedStemcells[stemcellConfig.Alias] != true {
				input.Stemcells = append(input.Stemcells, stemcellConfig)
			}
			usedStemcells[stemcellConfig.Alias] = true

			for j := range input.InstanceGroups {
				if input.InstanceGroups[j].VmType == vmType.Name {
					input.InstanceGroups[j].Stemcell = stemcellConfig.Alias
				}
			}
		}
	} else {
		for j := range input.InstanceGroups {
			input.InstanceGroups[j].Stemcell = ""
		}
	}

	return input
}
