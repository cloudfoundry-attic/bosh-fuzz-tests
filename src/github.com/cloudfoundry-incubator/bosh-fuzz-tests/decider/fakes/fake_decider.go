package fakes

type FakeDecider struct {
	IsYesYes bool
}

func (f *FakeDecider) IsYes() bool {
	return f.IsYesYes
}

type FakeDeciderMulti struct {
	YesResults []bool
	Default    bool
}

func (f *FakeDeciderMulti) IsYes() bool {
	var result bool
	if len(f.YesResults) > 0 {
		result, f.YesResults = f.YesResults[0], f.YesResults[1:]
	} else {
		result = f.Default
	}
	return result
}
