package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Parameter interface {
	Assign(input bftinput.Input) bftinput.Input
}
