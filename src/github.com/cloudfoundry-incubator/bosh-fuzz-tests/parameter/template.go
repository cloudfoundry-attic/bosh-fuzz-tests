package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type template struct {
	templates [][]string
}

func NewTemplate(templates [][]string) Parameter {
	return &template{
		templates: templates,
	}
}

func (t *template) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for j, _ := range input.InstanceGroups {
		pickedTemplates := t.templates[rand.Intn(len(t.templates))]
		input.InstanceGroups[j].Templates = []bftinput.Template{}

		for _, pickedTemplateName := range pickedTemplates {
			input.InstanceGroups[j].Templates = append(input.InstanceGroups[j].Templates, bftinput.Template{
				Name: pickedTemplateName,
			})
		}
	}

	return input
}
