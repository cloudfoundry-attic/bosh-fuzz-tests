package variables

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type PlaceholderPlanter interface {
	PlantPlaceholders(manifest *map[interface{}]interface{}, candidates [][]interface{}) (map[string]interface{}, error)
}

type placeholderPlanter struct {
	nameGenerator NameGenerator
}

func NewPlaceholderPlanter(nameGenerator NameGenerator) PlaceholderPlanter {
	return placeholderPlanter{
		nameGenerator: nameGenerator,
	}
}

func (p placeholderPlanter) PlantPlaceholders(manifest *map[interface{}]interface{}, candidates [][]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for _, candidate := range candidates {
		variableName := p.nameGenerator.Generate(8)
		record, err := p.setRecord(manifest, candidate, "(("+variableName+"))")
		if err != nil {
			return nil, err
		}
		result[variableName] = record
	}

	return result, nil
}

func (p placeholderPlanter) setRecord(manifest *map[interface{}]interface{}, path []interface{}, variableName string) (interface{}, error) {
	var result interface{}

	result = *manifest

	for index, item := range path {
		switch item.(type) {
		case string:
			mapItem := result.(map[interface{}]interface{})
			result = mapItem[item.(string)]

			if index == len(path)-1 {
				mapItem[item.(string)] = variableName
			}
		case int:
			arrayItem := result.([]interface{})
			result = arrayItem[item.(int)]

			if index == len(path)-1 {
				arrayItem[item.(int)] = variableName
			}
		}
	}

	return result, nil
}
