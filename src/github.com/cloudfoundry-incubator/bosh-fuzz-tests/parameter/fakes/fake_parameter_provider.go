package fakes

import (
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
)

type FakeParameterProvider struct {
	Stemcell *FakeStemcell
}

func NewFakeParameterProvider() *FakeParameterProvider {
	return &FakeParameterProvider{
		Stemcell: NewFakeStemcell(),
	}
}

func (p *FakeParameterProvider) Get(name string) bftparam.Parameter {
	if name == "stemcell" {
		return p.Stemcell
	}

	return nil
}
