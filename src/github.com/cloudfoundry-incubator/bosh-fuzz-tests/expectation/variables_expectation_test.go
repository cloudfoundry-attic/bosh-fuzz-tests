package expectation_test

import (
	"fmt"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VariablesExpectation", func() {
	var (
		variables   []input.Variable
		expectation Expectation
		cliRunner   *clirunnerfakes.FakeRunner
	)

	BeforeEach(func() {
		variables = createVariables(3)
		expectation = NewVariablesExpectation(variables)
		cliRunner = &clirunnerfakes.FakeRunner{}
	})

	Context("when expected variables match the number of variables created", func() {
		BeforeEach(func() {
			cliRunner.RunWithOutputReturns(EventLogWithEvents, nil)
		})

		It("does not return an error", func() {
			err := expectation.Run(cliRunner, "1")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when expected variables does NOT match the number of variables created", func() {
		BeforeEach(func() {
			cliRunner.RunWithOutputReturns(EventLogWithNoEvents, nil)
		})

		It("returns an error", func() {
			err := expectation.Run(cliRunner, "1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Expected %d variables to be created but found 0", len(variables)))
		})
	})
})

func createVariables(count int) []input.Variable {
	result := []input.Variable{}
	for i := 0; i < count; i++ {
		variable := input.Variable{Name: fmt.Sprintf("/TestDirector/foo-deployment/%d", i)}
		result = append(result, variable)
	}
	return result
}

var EventLogWithNoEvents = `
{
    "Tables": [
        {
            "Rows": []
        }
    ]
}`

var EventLogWithEvents = `
{
    "Tables": [
        {
            "Rows": [
                {
                    "action": "create",
                    "context": "id: \"19\"\nname: /TestDirector/foo-deployment/0",
                    "deployment": "foo-deployment",
                    "error": "",
                    "id": "71",
                    "instance": "",
                    "object_name": "/TestDirector/foo-deployment/0",
                    "object_type": "variable",
                    "task_id": "5",
                    "time": "Mon May  8 15:17:24 UTC 2017",
                    "user": "test"
                },
                {
                    "action": "create",
                    "context": "id: \"18\"\nname: /TestDirector/foo-deployment/1",
                    "deployment": "foo-deployment",
                    "error": "",
                    "id": "70",
                    "instance": "",
                    "object_name": "/TestDirector/foo-deployment/1",
                    "object_type": "variable",
                    "task_id": "5",
                    "time": "Mon May  8 15:17:24 UTC 2017",
                    "user": "test"
                },
                {
                    "action": "create",
                    "context": "id: \"17\"\nname: /TestDirector/foo-deployment/2",
                    "deployment": "foo-deployment",
                    "error": "",
                    "id": "69",
                    "instance": "",
                    "object_name": "/TestDirector/foo-deployment/2",
                    "object_type": "variable",
                    "task_id": "5",
                    "time": "Mon May  8 15:17:23 UTC 2017",
                    "user": "test"
                }
            ]
        }
    ]
}`
