package variables

import (
	"reflect"
)

type PathPicker interface {
	PickPaths(paths [][]interface{}, requestedPicks int) [][]interface{}
}

type pathPicker struct {
	randomizer NumberRandomizer
}

func NewPathPicker(randomizer NumberRandomizer) PathPicker {
	return pathPicker{
		randomizer: randomizer,
	}
}

// number of picks returned is not guaranteed to match requestedPicks
func (p pathPicker) PickPaths(paths [][]interface{}, requestedPicks int) [][]interface{} {
	picks := [][]interface{}{}

	for i := 0; len(paths) > 0 && i < requestedPicks; i++ {
		index := p.randomizer.Intn(len(paths))
		value := paths[index]

		paths = p.trimmedPaths(paths, value)

		picks = append(picks, value)
	}

	return picks
}

func (p pathPicker) compareStartsWithPath(a []interface{}, b []interface{}) bool {
	shorterPath := a
	longerPath := b
	if len(b) < len(a) {
		shorterPath = b
		longerPath = a
	}

	for index, shorterPathElement := range shorterPath {
		longerPathElement := longerPath[index]
		if reflect.TypeOf(shorterPathElement) != reflect.TypeOf(longerPathElement) || shorterPathElement != longerPathElement {
			return false
		}
	}
	return true
}

func (p pathPicker) trimmedPaths(paths [][]interface{}, toRemove []interface{}) [][]interface{} {
	result := [][]interface{}{}

	for _, value := range paths {
		if !p.compareStartsWithPath(value, toRemove) {
			result = append(result, value)
		}
	}

	return result
}
