package parameter_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {
	var (
		stemcell Parameter
	)

	Context("when definition is os", func() {
		BeforeEach(func() {
			stemcell = NewStemcell("os")
		})

		Context("when input has vm types", func() {
			input := 
		})

		Context("when input has resource pools", func() {

		})
	})

	Context("when definition is name", func() {
		BeforeEach(func() {
			stemcell = NewStemcell("name")
		})
	})
})
