package deployment_test

import (
	"errors"

	"strings"

	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer/analyzerfakes"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment/deploymentfakes"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation/expectationfakes"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables/variablesfakes"
	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployer", func() {
	var (
		cliRunner       *clirunnerfakes.FakeRunner
		uaaRunner       *clirunnerfakes.FakeRunner
		renderer        *deploymentfakes.FakeRenderer
		inputGenerator  *deploymentfakes.FakeInputGenerator
		analyzer        *analyzerfakes.FakeAnalyzer
		fs              *fakesys.FakeFileSystem
		sprinkler       *variablesfakes.FakeSprinkler
		errandGenerator *deploymentfakes.FakeStepGenerator

		deployer Deployer
	)

	BeforeEach(func() {
		cliRunner = &clirunnerfakes.FakeRunner{}
		uaaRunner = &clirunnerfakes.FakeRunner{}
		renderer = &deploymentfakes.FakeRenderer{}
		inputGenerator = &deploymentfakes.FakeInputGenerator{}
		analyzer = &analyzerfakes.FakeAnalyzer{}
		fs = fakesys.NewFakeFileSystem()
		sprinkler = &variablesfakes.FakeSprinkler{}
		errandGenerator = &deploymentfakes.FakeStepGenerator{}

		directorInfo := bltaction.DirectorInfo{
			Name: "fake-director",
			UUID: "fake-director-uuid",
			URL:  "fake-director-url",
		}

		logger := boshlog.NewLogger(boshlog.LevelNone)

		deployer = NewDeployer(cliRunner, uaaRunner, directorInfo, renderer, inputGenerator, []StepGenerator{errandGenerator}, analyzer, sprinkler, fs, logger, false)
	})

	Context("when fs errors when creating temporary file", func() {
		BeforeEach(func() {
			fs.TempFileError = errors.New("error")
		})

		It("should also return an error", func() {
			err := deployer.RunDeploys()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when input generator returns an error", func() {
		BeforeEach(func() {
			inputGenerator.GenerateReturns(nil, errors.New("error"))
		})

		It("should also return an error", func() {
			err := deployer.RunDeploys()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Generating input: error"))
		})
	})

	Context("when analyzer has test cases", func() {
		var cases []bftanalyzer.Case

		BeforeEach(func() {
			cases = []bftanalyzer.Case{{}}
			analyzer.AnalyzeReturns(cases)
		})

		Context("when renderer fails", func() {
			BeforeEach(func() {
				renderer.RenderReturns(errors.New("error"))
			})

			It("returns an error", func() {
				err := deployer.RunDeploys()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Rendering deployment manifest: error"))
			})
		})

		Context("when trying to upload cloud-config", func() {
			Context("when cli runner fails", func() {
				BeforeEach(func() {
					cliRunner.RunWithArgsStub = func(args ...string) error {
						if args[0] == "update-cloud-config" {
							return errors.New("error")
						}
						return nil
					}
				})

				It("returns an error", func() {
					err := deployer.RunDeploys()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Updating cloud config: error"))
				})
			})
		})

		Context("when trying to deploy", func() {
			Context("when cli runner fails", func() {
				BeforeEach(func() {
					cliRunner.RunWithOutputStub = func(args ...string) (string, error) {
						if strings.Join(args[:3], " ") == "-d foo-deployment deploy" {
							return "Task 1", errors.New("error")
						}
						return "", nil
					}
				})

				It("returns an error", func() {
					err := deployer.RunDeploys()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Running deploy: error"))
				})

				Context("when deployment failure is expected", func() {
					BeforeEach(func() {
						cases[0].DeploymentWillFail = true
					})

					It("does not returns an error", func() {
						err := deployer.RunDeploys()
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
		})

		Context("when the test cases have expectations", func() {
			BeforeEach(func() {
				cliRunner.RunWithOutputReturns("Task 1", nil)
			})

			Context("when expecatation fails", func() {
				BeforeEach(func() {
					fakeExpectation := &expectationfakes.FakeExpectation{}
					fakeExpectation.RunReturns(errors.New("error"))
					cases[0].Expectations = append(cases[0].Expectations, fakeExpectation)
				})

				It("returns an error", func() {
					err := deployer.RunDeploys()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Running expectation: error"))
				})
			})
		})

		Context("when there are errand steps", func() {
			BeforeEach(func() {
				cliRunner.RunWithOutputReturns("Task 1", nil)
			})

			var fakeStep *deploymentfakes.FakeStep

			BeforeEach(func() {
				fakeStep = &deploymentfakes.FakeStep{}
				errandGenerator.StepsReturns([]Step{fakeStep})
			})

			It("runs the steps", func() {
				err := deployer.RunDeploys()
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeStep.RunCallCount()).To(Equal(1))
				Expect(fakeStep.RunArgsForCall(0)).To(Equal(cliRunner))
			})
		})
	})
})
