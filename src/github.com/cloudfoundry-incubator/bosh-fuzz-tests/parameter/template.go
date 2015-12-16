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
	for j, _ := range input.Jobs {
		pickedTemplates := t.templates[rand.Intn(len(t.templates))]
		input.Jobs[j].Templates = []bftinput.Template{}

		for _, pickedTemplateName := range pickedTemplates {
			input.Jobs[j].Templates = append(input.Jobs[j].Templates, bftinput.Template{
				Name: pickedTemplateName,
			})
		}
	}

	return input
}
