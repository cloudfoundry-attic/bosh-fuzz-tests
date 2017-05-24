package variables


type PathWeeder interface {
	WeedPaths(paths [][]interface{}, invalidPaths [][]interface{}) [][]interface{}
}

type pathWeeder struct {
}


type Anything struct {

}

var Integer Anything = Anything{}
var String Anything = Anything{}


func NewPathWeeder() PathWeeder {
	return pathWeeder{
	}
}

var badPathPatterns [][]interface{} = [][]interface{} {
	{"properties"},
	{"instance_groups", Integer, "properties"},
	{"instance_groups", Integer, "jobs", Integer, "properties"},
	{"instance_groups", Integer, "jobs", Integer, "consumes", String, "properties"},
	{"jobs", Integer, "properties"},
	{"jobs", Integer, "templates", Integer, "properties"},
	{"jobs", Integer, "templates", Integer, "consumes", String, "properties"},
	{"instance_groups", Integer, "env"},
	{"jobs", Integer, "env"},
	{"resource_pools", Integer, "env"},
}

// Returns:
// array of locations of possible placeholder locations in the format:
// map item name or array index
// for example [ ['hi', 3, 'there', 'property'], ]
// and invalid ones are [ ['hi', Anything, 'there', Anything], ]

func (b pathWeeder) WeedPaths(paths [][]interface{}, invalidPaths [][]interface{}) [][]interface{} {
	return paths
}

func pathMatchesPattern(path []interface{}, pattern []interface{}) bool {
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
