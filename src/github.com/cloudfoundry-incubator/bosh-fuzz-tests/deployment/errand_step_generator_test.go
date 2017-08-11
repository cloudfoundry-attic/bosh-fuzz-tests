package deployment_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrandStepGenerator", func() {
	var (
		generator     deployment.ErrandStepGenerator
		testCase      analyzer.Case
		testJobs      []bftinput.Job
		testTemplates []bftinput.Template
	)

	BeforeEach(func() {
		testJobs = nil
		generator = deployment.NewErrandStepGenerator()
	})

	Describe("Steps", func() {
		JustBeforeEach(func() {
			if testJobs == nil {
				testJobs = []bftinput.Job{{Name: "instance-name", Templates: testTemplates}}
			}
			testCase = analyzer.Case{
				Input: bftinput.Input{
					Jobs: testJobs,
				},
			}
		})

		BeforeEach(func() {
			testTemplates = []bftinput.Template{{Name: "template-name"}}
		})

		It("returns an errand step that has the correct name and deployment name", func() {
			Expect(generator.Steps(testCase)).To(Equal([]deployment.Step{deployment.ErrandStep{Name: "template-name", DeploymentName: "foo-deployment"}}))
		})

		Context("when input's job has no templates", func() {
			BeforeEach(func() {
				testTemplates = []bftinput.Template{}
			})

			It("returns an empty array of Steps", func() {
				Expect(generator.Steps(testCase)).To(Equal([]deployment.Step{}))
			})
		})

		Context("when input has no jobs", func() {
			BeforeEach(func() {
				testJobs = []bftinput.Job{}
			})

			It("returns an empty array of Steps", func() {
				Expect(generator.Steps(testCase)).To(Equal([]deployment.Step{}))
			})
		})
	})
})

var _ = Describe("ErrandStep", func() {
	var (
		cliRunner *clirunnerfakes.FakeRunner
		step      deployment.ErrandStep
	)

	BeforeEach(func() {
		cliRunner = &clirunnerfakes.FakeRunner{}
		step = deployment.ErrandStep{Name: "yogurt", DeploymentName: "greek"}
	})

	Describe("Run", func() {
		It("runs an errand command", func() {
			err := step.Run(cliRunner)
			Expect(err).NotTo(HaveOccurred())

			Expect(cliRunner.RunWithArgsCallCount()).To(Equal(1))
			args := cliRunner.RunWithArgsArgsForCall(0)
			Expect(args).To(HaveLen(4))
			Expect(args[0]).To(Equal("run-errand"))
			Expect(args[1]).To(Equal("yogurt"))
			Expect(args[2]).To(Equal("-d"))
			Expect(args[3]).To(Equal("greek"))
		})

		Context("when cli runner errors", func() {
			BeforeEach(func() {
				cliRunner.RunWithArgsReturns(errors.New("I'm an error"))
			})

			It("bubbles the error up", func() {
				err := step.Run(cliRunner)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("I'm an error"))
			})
		})
	})
})
