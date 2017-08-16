package deployment_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrandStepGenerator", func() {
	var (
		generator          deployment.ErrandStepGenerator
		testCase           analyzer.Case
		testInstanceGroups []bftinput.InstanceGroup
	)

	BeforeEach(func() {
		generator = deployment.NewErrandStepGenerator()
	})

	Describe("Steps", func() {
		JustBeforeEach(func() {
			testCase = analyzer.Case{
				Input: bftinput.Input{
					InstanceGroups: testInstanceGroups,
				},
			}
		})

		BeforeEach(func() {
			testJobs := []bftinput.Job{
				{Name: "job-name"},
				{Name: "other-job-name"},
			}
			secondInstanceGroupTestJobs := []bftinput.Job{
				{Name: "other-instance-group-job-name"},
			}
			testInstanceGroups = []bftinput.InstanceGroup{
				{Name: "instance-name", Jobs: testJobs, Instances: 1},
				{Name: "other-instance-group", Jobs: secondInstanceGroupTestJobs, Instances: 1},
			}
		})

		DescribeTable("instance filters", func(name, instanceFilter string) {
			Eventually(func() []deployment.Step {
				return generator.Steps(testCase)
			}, time.Second, time.Microsecond).Should(ContainElement(
				deployment.ErrandStep{
					Name:             name,
					InstanceFilter:   instanceFilter,
					DeploymentName:   "foo-deployment",
					CommandLineFlags: []string{},
				},
			))
		},
			Entry("", "job-name", ""),
			Entry("", "job-name", "instance-name"),
			Entry("", "job-name", "instance-name/0"),

			Entry("", "other-job-name", ""),
			Entry("", "other-job-name", "instance-name"),
			Entry("", "other-job-name", "instance-name/0"),

			Entry("", "other-instance-group-job-name", ""),
			Entry("", "other-instance-group-job-name", "other-instance-group"),
			Entry("", "other-instance-group-job-name", "other-instance-group/0"),
		)

		DescribeTable("command line options", func(flags []string) {
			Eventually(func() []deployment.Step {
				return generator.Steps(testCase)
			}, time.Second, time.Microsecond).Should(ContainElement(
				deployment.ErrandStep{
					Name:             "job-name",
					DeploymentName:   "foo-deployment",
					CommandLineFlags: flags,
				},
			))
		},
			Entry("can be empty", []string{}),
			Entry("can have keep-alive", []string{"keep-alive"}),
		)

		DescribeTable("number of steps returned", func(numberOfSteps int) {
			Eventually(func() []deployment.Step {
				return generator.Steps(testCase)
			}, time.Second, time.Microsecond).Should(HaveLen(numberOfSteps))
		},
			Entry("", 0),
			Entry("", 1),
			Entry("", 2),
			Entry("", 3),
			Entry("", 4),
		)

		Context("when instance group has lifecycle errand", func() {
			BeforeEach(func() {
				testInstanceGroups = []bftinput.InstanceGroup{{
					Name:      "instance-name",
					Jobs:      []bftinput.Job{{Name: "job-name"}},
					Instances: 1,
					Lifecycle: "errand",
				}}
			})

			DescribeTable("name of errand step", func(name string) {
				Eventually(func() []deployment.Step {
					return generator.Steps(testCase)
				}, time.Second, time.Microsecond).Should(ContainElement(
					deployment.ErrandStep{
						Name:             name,
						DeploymentName:   "foo-deployment",
						CommandLineFlags: []string{},
					},
				))
			},
				Entry("is sometimes instance group", "instance-name"),
				Entry("is sometimes job name", "job-name"),
			)
		})

		Context("when instance group has lifecycle service", func() {
			BeforeEach(func() {
				testInstanceGroups = []bftinput.InstanceGroup{{
					Name:      "instance-name",
					Jobs:      []bftinput.Job{{Name: "job-name"}},
					Instances: 1,
					Lifecycle: "service",
				}}
			})

			It("should never set errand step name to the instance group name", func() {
				Consistently(func() []deployment.Step {
					return generator.Steps(testCase)
				}, 50*time.Millisecond, time.Microsecond).ShouldNot(ContainElement(
					deployment.ErrandStep{
						Name:             "instance-name",
						DeploymentName:   "foo-deployment",
						CommandLineFlags: []string{},
					},
				))
			})
		})

		Context("when input's instance group has no jobs", func() {
			BeforeEach(func() {
				testInstanceGroups = []bftinput.InstanceGroup{{Name: "instance-name", Jobs: []bftinput.Job{}, Instances: 1}}
			})

			It("returns an empty array of Steps", func() {
				Expect(generator.Steps(testCase)).To(Equal([]deployment.Step{}))
			})
		})

		Context("when input has no instance groups", func() {
			BeforeEach(func() {
				testInstanceGroups = []bftinput.InstanceGroup{}
			})

			It("returns an empty array of Steps", func() {
				Expect(generator.Steps(testCase)).To(Equal([]deployment.Step{}))
			})
		})

		Context("when input's instance group has no instances", func() {
			BeforeEach(func() {
				testInstanceGroups = []bftinput.InstanceGroup{{Name: "instance-name", Jobs: []bftinput.Job{{Name: "job"}}, Instances: 0}}
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

		Context("when cli flags are present", func() {
			BeforeEach(func() {
				step = deployment.ErrandStep{
					Name:             "yogurt",
					DeploymentName:   "plain",
					CommandLineFlags: []string{"fancy", "behavior"},
				}
			})

			It("sends those flags to the cli", func() {
				err := step.Run(cliRunner)
				Expect(err).NotTo(HaveOccurred())

				Expect(cliRunner.RunWithArgsCallCount()).To(Equal(1))
				args := cliRunner.RunWithArgsArgsForCall(0)
				Expect(args).To(HaveLen(6))
				Expect(args[0]).To(Equal("run-errand"))
				Expect(args[1]).To(Equal("yogurt"))
				Expect(args[2]).To(Equal("-d"))
				Expect(args[3]).To(Equal("plain"))
				Expect(args[4]).To(Equal("--fancy"))
				Expect(args[5]).To(Equal("--behavior"))
			})
		})

		Context("when an instance filter is present", func() {
			BeforeEach(func() {
				step = deployment.ErrandStep{
					Name:           "yogurt",
					DeploymentName: "plain",
					InstanceFilter: "fruitatthebottom",
				}
			})

			It("filters the errand with --instance flag", func() {
				err := step.Run(cliRunner)
				Expect(err).NotTo(HaveOccurred())

				Expect(cliRunner.RunWithArgsCallCount()).To(Equal(1))
				args := cliRunner.RunWithArgsArgsForCall(0)
				Expect(args).To(HaveLen(6))
				Expect(args[0]).To(Equal("run-errand"))
				Expect(args[1]).To(Equal("yogurt"))
				Expect(args[2]).To(Equal("-d"))
				Expect(args[3]).To(Equal("plain"))
				Expect(args[4]).To(Equal("--instance"))
				Expect(args[5]).To(Equal("fruitatthebottom"))
			})
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
