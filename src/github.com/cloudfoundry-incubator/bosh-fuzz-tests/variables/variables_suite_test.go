package variables_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVariables(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Variables Suite")
}
