package variables

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"gopkg.in/yaml.v2"
)

type Sprinkler interface {
	SprinklePlaceholders(manifestPath string, badPathFilter [][]interface{}) (map[string]interface{}, error)
}

type sprinkler struct {
	parameters         bftconfig.Parameters
	fs                 boshsys.FileSystem
	randomizer         NumberRandomizer
	pathBuilder        PathBuilder
	pathPicker         PathPicker
	placeholderPlanter PlaceholderPlanter
}

func NewSprinkler(
	parameters bftconfig.Parameters,
	fs boshsys.FileSystem,
	randomizer NumberRandomizer,
	pathBuilder PathBuilder,
	pathPicker PathPicker,
	placeholderPlanter PlaceholderPlanter,
) Sprinkler {
	return &sprinkler{
		parameters:         parameters,
		fs:                 fs,
		randomizer:         randomizer,
		pathBuilder:        pathBuilder,
		pathPicker:         pathPicker,
		placeholderPlanter: placeholderPlanter,
	}
}

// Returns: map
// key is the placeholder name. value is the placeholder value
func (s sprinkler) SprinklePlaceholders(manifestPath string, badPathFilter [][]interface{}) (map[string]interface{}, error) {
	manifest := map[interface{}]interface{}{}

	yamlFile, err := s.fs.ReadFile(manifestPath)
	if err != nil {
		return nil, bosherr.WrapError(err, "Error reading manifest file")
	}

	err = yaml.Unmarshal(yamlFile, manifest)
	if err != nil {
		return nil, bosherr.WrapError(err, "Error unmarshalling manifest file")
	}

	placeholderPaths := s.pathBuilder.BuildPaths(manifest)
	randomizer := DefaultNumberRandomizer{}

	placeholderPaths = NewPathWeeder().WeedPaths(placeholderPaths, badPathFilter)

	numOfSubstitutions := s.parameters.NumOfSubstitutions[randomizer.Intn(len(s.parameters.NumOfSubstitutions))]
	candidates := s.pathPicker.PickPaths(placeholderPaths, numOfSubstitutions)

	substitutedVariables, err := s.placeholderPlanter.PlantPlaceholders(&manifest, candidates) //-- returns a list of values that were substituted
	if err != nil {
		return nil, bosherr.WrapError(err, "Error adding variables to manifest file")
	}

	manifestString, err := yaml.Marshal(manifest)
	if err != nil {
		return nil, bosherr.WrapError(err, "Error marshalling manifest file")
	}

	err = s.fs.WriteFile(manifestPath, manifestString)
	if err != nil {
		return nil, bosherr.WrapError(err, "Error writing manifest file")
	}

	return substitutedVariables, nil
}
