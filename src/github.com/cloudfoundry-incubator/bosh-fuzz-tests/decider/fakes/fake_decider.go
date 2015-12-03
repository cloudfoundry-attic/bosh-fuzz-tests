package fakes

type FakeDecider struct {
	IsYesYes bool
}

func (f *FakeDecider) IsYes() bool {
	return f.IsYesYes
}
