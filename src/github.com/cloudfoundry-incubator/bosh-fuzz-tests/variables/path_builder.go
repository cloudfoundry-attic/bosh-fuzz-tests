package variables

type PathBuilder interface {
	BuildPaths(obj interface{}) [][]interface{}
}

type pathBuilder struct {
}

func NewPathBuilder() PathBuilder {
	return pathBuilder{}
}

// Returns:
// array of locations of possible placeholder locations in the format:
// map item name or array index
// for example [ ['hi', 3, 'there', 'property'], ]
// for example [ ['jobs', 'templates', 0, 'name'], ]
func (b pathBuilder) BuildPaths(obj interface{}) [][]interface{} {
	output := [][]interface{}{}
	b.buildPaths(nil, obj, &output)

	return output
}

func (b pathBuilder) buildPaths(path []interface{}, obj interface{}, output *[][]interface{}) {
	if path == nil {
		path = []interface{}{}
	}

	switch obj.(type) {
	case []interface{}:
		for index, item := range obj.([]interface{}) {
			newPath := b.appendPath(path, index)
			*output = append(*output, newPath)
			b.buildPaths(newPath, item, output)
		}
	case map[interface{}]interface{}:
		for key, value := range obj.(map[interface{}]interface{}) {
			newPath := b.appendPath(path, key)
			*output = append(*output, newPath)
			b.buildPaths(newPath, value, output)
		}
	}
}

func (b pathBuilder) appendPath(path []interface{}, value interface{}) []interface{} {
	newPath := []interface{}{}
	newPath = append(newPath, path...)
	newPath = append(newPath, value)

	return newPath
}
