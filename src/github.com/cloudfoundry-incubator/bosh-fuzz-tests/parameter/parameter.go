package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Parameter interface {
	Apply(input bftinput.Input) bftinput.Input
}
