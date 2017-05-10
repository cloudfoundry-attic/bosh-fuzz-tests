package input

import "reflect"

type Variable struct {
	Name    string
	Type    string
	Options map[string]interface{}
}

func (v Variable) IsEqual(other Variable) bool {
	return reflect.DeepEqual(v, other)
}
