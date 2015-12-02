package fakes

type FakeNameGenerator struct {
	Names []string
}

func (f *FakeNameGenerator) Generate(length int) string {
	var name string
	name, f.Names = f.Names[0], f.Names[1:]
	return name
}
