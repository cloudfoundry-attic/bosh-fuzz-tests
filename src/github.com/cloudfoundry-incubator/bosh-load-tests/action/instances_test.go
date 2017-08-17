package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/bosh-load-tests/action"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
	"errors"
)

var _ = Describe("GetInstances", func() {
	var (
		cliRunner      *clirunnerfakes.FakeRunner
		directorInfo   action.DirectorInfo
		deploymentName string
		instancesInfo  *action.InstancesInfo
	)

	BeforeEach(func() {
		cliRunner = &clirunnerfakes.FakeRunner{}
		directorInfo = action.DirectorInfo{
			UUID: "director-uuid",
			URL:  "https://example.com",
			Name: "My Little Director",
		}
		deploymentName = "my-deployment"

		instancesInfo = action.NewInstances(directorInfo, deploymentName, cliRunner)
	})

	It("runs instances command for given deployment and parses out IDs of instances", func() {
		cliRunner.RunWithOutputReturns(`{
    "Tables": [
        {
            "Content": "instances",
            "Header": {
                "az": "AZ",
                "instance": "Instance",
                "ips": "IPs",
                "process_state": "Process State"
            },
            "Rows": [
                {
                    "az": "z1",
                    "instance": "smoke-tests/5b24e1b6-2a19-44ad-b966-7bbbc06ef7b0",
                    "ips": "",
                    "process_state": ""
                },
                {
                    "az": "z2",
                    "instance": "zookeeper/35cadb81-95be-441d-838e-898f2cc8a4a1",
                    "ips": "10.244.0.4",
                    "process_state": "running"
                }
            ],
            "Notes": null
        }
    ],
    "Blocks": null,
    "Lines": [
        "Using environment '192.168.50.6' as client 'admin'",
        "Task 4",
        ". Done",
        "Succeeded"
    ]
}`, nil)

		foundInstances, err := instancesInfo.GetInstances()
		Expect(err).ToNot(HaveOccurred())
		Expect(cliRunner.RunWithOutputCallCount()).To(Equal(1))
		Expect(cliRunner.RunWithOutputArgsForCall(0)).To(Equal([]string{
			"-d",
			"my-deployment",
			"instances",
			"--json",
		}))

		Expect(foundInstances).To(Equal(map[string][]action.Instance{
			"smoke-tests": {{Name: "smoke-tests", ID: "5b24e1b6-2a19-44ad-b966-7bbbc06ef7b0"}},
			"zookeeper":   {{Name: "zookeeper", ID: "35cadb81-95be-441d-838e-898f2cc8a4a1"}},
		}))
	})

	Context("When CLI command fails", func() {
		It("Returns an empty map and error from cli", func() {
			cliRunner.RunWithOutputReturns("", errors.New("CLI failed"))

			foundInstances, err := instancesInfo.GetInstances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("CLI failed"))
			Expect(foundInstances).To(Equal(map[string][]action.Instance{}))
		})
	})

	Context("When CLI returns invalid JSON", func() {
		It("Returns an empty and error from json unmarshalling", func() {
			cliRunner.RunWithOutputReturns(`foo::: monkey`, nil)

			foundInstances, err := instancesInfo.GetInstances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error unmarshalling JSON"))
			Expect(foundInstances).To(Equal(map[string][]action.Instance{}))
		})
	})

	Context("when output contains empty tables record", func() {
		It("returns empty map of instances", func() {
			cliRunner.RunWithOutputReturns(`{"Tables": []}`, nil)

			foundInstances, err := instancesInfo.GetInstances()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundInstances).To(Equal(map[string][]action.Instance{}))
		})
	})

	Context("when output contains empty rows in table", func() {
		It("returns empty map of instances", func() {
			cliRunner.RunWithOutputReturns(`{"Tables": [{"Rows":[]}]}`, nil)

			foundInstances, err := instancesInfo.GetInstances()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundInstances).To(Equal(map[string][]action.Instance{}))
		})
	})

	Context("when output contains rows without 'instance' column in table", func() {
		It("returns empty map of instances", func() {
			cliRunner.RunWithOutputReturns(`{"Tables": [{"Rows":[{"foo":"bar"}]}]}`, nil)

			foundInstances, err := instancesInfo.GetInstances()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundInstances).To(Equal(map[string][]action.Instance{}))
		})
	})

	Context("when output contains rows with 'instance' column without slash in table", func() {
		It("returns empty map of instances", func() {
			cliRunner.RunWithOutputReturns(`{"Tables": [{"Rows":[{"instance":"my old instance"}]}]}`, nil)

			foundInstances, err := instancesInfo.GetInstances()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundInstances).To(Equal(map[string][]action.Instance{}))
		})
	})
})
