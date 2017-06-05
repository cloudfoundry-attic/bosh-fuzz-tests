package variables_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PathWeeder", func() {

	Context("from a list of paths", func() {
		var paths, badCloudConfigPatterns [][]interface{}
		BeforeEach(func() {
			paths = [][]interface{}{
				{"azs", 0, "iaas"},
				{"director_name"},
				{"instance_groups", 0, "jobs", 0, "properties"},
				{"instance_groups", 0, "jobs", 1, "consumes", "dblinks", "properties"},
				{"instance_groups", 0, "properties"},
				{"name"},
				{"network", "test", "jobs", "integer"},
				{"network", "test", "jobs"},
				{"network", "test", "properties"},
				{"properties", "autocorrect"},
				{"properties"},
				{"releases", 0, "test-release-name"},
				{"releases", 0},
				{"releases"},
				{"stemcells", 0},
				{"stemcells"},
			}
			badCloudConfigPatterns = [][]interface{}{
				{"resource_pools"},
				{"instance_groups", Integer, "jobs", Integer, "properties"},
				{"instance_groups", Integer, "jobs", Integer, String, String, "properties"},
				{"instance_groups", Integer, "properties"},
				{"name"},
				{"properties"},
				{"releases", Integer, String},
				{"releases", Integer},
				{"releases"},
				{"stemcells", Integer},
				{"stemcells"},
			}
		})

		It("has all invalid paths removed", func() {
			var expectedPaths [][]interface{} = [][]interface{}{
				{"azs", 0, "iaas"},
				{"director_name"},
				{"network", "test", "jobs", "integer"},
				{"network", "test", "jobs"},
				{"network", "test", "properties"},
				{"properties", "autocorrect"},
			}
			result := NewPathWeeder().WeedPaths(paths, badCloudConfigPatterns)
			Expect(result).To(Equal(expectedPaths))
		})
	})

})
