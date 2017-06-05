package variables

type PathWeeder interface {
	WeedPaths(paths [][]interface{}, badPathPatterns [][]interface{}) [][]interface{}
}

type pathWeeder struct {
}

type Anything struct {
}

var Integer Anything = Anything{}
var String Anything = Anything{}

func NewPathWeeder() PathWeeder {
	return pathWeeder{}
}

// Returns:
// array of locations of possible placeholder locations in the format:
// map item name or array index
// for example [ ['hi', 3, 'there', 'property'], ]
// and invalid ones are [ ['hi', Anything, 'there', Anything], ]

func (p pathWeeder) WeedPaths(paths [][]interface{}, badPathPatterns [][]interface{}) [][]interface{} {

	for _, pattern := range badPathPatterns {
		paths = p.trimmedPaths(paths, pattern)
	}
	return paths
}

func (p pathWeeder) trimmedPaths(paths [][]interface{}, patternToRemove []interface{}) [][]interface{} {
	result := [][]interface{}{}

	for _, value := range paths {
		if !p.pathMatchesPattern(value, patternToRemove) {
			result = append(result, value)
		}
	}

	return result
}

func (p pathWeeder) pathMatchesPattern(path []interface{}, pattern []interface{}) bool {
	if len(path) != len(pattern) {
		return false
	}

	for index, value := range pattern {
		switch value.(type) {
		case string:
			if path[index].(string) != value.(string) {
				return false
			}
		case int:
			if path[index].(int) != value.(int) {
				return false
			}
		}
	}

	return true
}
