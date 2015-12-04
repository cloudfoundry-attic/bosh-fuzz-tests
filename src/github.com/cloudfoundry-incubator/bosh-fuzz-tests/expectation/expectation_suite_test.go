package expectation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestExpectation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Expectation Suite")
}
