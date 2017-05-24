package variables_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("PathBuilder", func() {

	var yamlString string
	var yamlMap map[interface{}]interface{}
	var pathBuilder PathBuilder

	BeforeEach(func() {
		pathBuilder = NewPathBuilder()
	})

	Context("when user gives an empty path", func() {
		BeforeEach(func() {
			yamlString = `
name: foo-deployment
`
			yamlMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(yamlString), yamlMap)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should still work", func() {
			paths := pathBuilder.BuildPaths(yamlMap)
			Expect(paths).ToNot(BeNil())
			Expect(len(paths)).To(Equal(1))
			Expect(paths).To(Equal([][]interface{}{{"name"}}))
		})
	})

	Context("when yaml contains simple elements", func() {
		BeforeEach(func() {
			yamlString = `
version: nothing
release: fake
name: foo-deployment
`
			yamlMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(yamlString), yamlMap)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return the correct paths", func() {
			paths := pathBuilder.BuildPaths(yamlMap)
			Expect(paths).ToNot(BeNil())
			Expect(len(paths)).To(Equal(3))
			Expect(paths).To(ConsistOf([][]interface{}{{"name"}, {"version"}, {"release"}}))
		})
	})

	Context("when yaml contains nested elements", func() {
		BeforeEach(func() {
			yamlString = `
name:
  version:
    major: 1
    minor: 2
    patch: 3
  release: fake
`
			yamlMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(yamlString), yamlMap)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return the correct paths", func() {
			paths := pathBuilder.BuildPaths(yamlMap)
			Expect(paths).ToNot(BeNil())
			Expect(len(paths)).To(Equal(6))
			Expect(paths).To(ConsistOf([][]interface{}{
				{"name"},
				{"name", "version"},
				{"name", "version", "major"},
				{"name", "version", "minor"},
				{"name", "version", "patch"},
				{"name", "release"},
			}))
		})
	})

	Context("when yaml contains arrays", func() {
		BeforeEach(func() {
			yamlString = `
names: [a, b, c, d]
ranks:
- direction: west
  altitude: 50
- direction: east
  altitude: 25
- direction: south
  altitude: 12
`
			yamlMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(yamlString), yamlMap)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return the correct paths", func() {
			paths := pathBuilder.BuildPaths(yamlMap)
			Expect(paths).ToNot(BeNil())
			Expect(len(paths)).To(Equal(15))
			Expect(paths).To(ConsistOf([][]interface{}{
				{"names"},
				{"names", 0},
				{"names", 1},
				{"names", 2},
				{"names", 3},
				{"ranks"},
				{"ranks", 0},
				{"ranks", 0, "direction"},
				{"ranks", 0, "altitude"},
				{"ranks", 1},
				{"ranks", 1, "direction"},
				{"ranks", 1, "altitude"},
				{"ranks", 2},
				{"ranks", 2, "direction"},
				{"ranks", 2, "altitude"},
			}))
		})
	})

	Context("when yaml contains nested and array elements",func() {
		BeforeEach(func() {
			yamlString = `
officer:
  ranks:
  - direction: west
    altitude: 50
  - direction: east
    altitude: 25
  - direction: south
    altitude: 12
`
			yamlMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(yamlString), yamlMap)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return the correct paths", func() {
			paths := pathBuilder.BuildPaths(yamlMap)
			Expect(paths).ToNot(BeNil())
			Expect(len(paths)).To(Equal(11))
			Expect(paths).To(ConsistOf([][]interface{}{
				{"officer"},
				{"officer", "ranks"},
				{"officer", "ranks", 0},
				{"officer", "ranks", 0, "direction"},
				{"officer", "ranks", 0, "altitude"},
				{"officer", "ranks", 1},
				{"officer", "ranks", 1, "direction"},
				{"officer", "ranks", 1, "altitude"},
				{"officer", "ranks", 2},
				{"officer", "ranks", 2, "direction"},
				{"officer", "ranks", 2, "altitude"},
			}))
		})
	})

	Context("when yaml contains simple, nested and array elements",func() {
		BeforeEach(func() {
			yamlString = `
name: foo-deployment

releases:
- name: foo-release
  version: 0+dev.1
- name: bar-release
  version: 30+dev.6

stemcells:
- alias: stemcell-2
  name: ubuntu-stemcell
  version: "latest"

update:
  canaries: 2
  canary_watch_time: 4000
  max_in_flight: 2
  update_watch_time: 20

jobs:
  name: zRD
  instances: 5
  persistent_disk_type: czcuBXB7WY
  stemcell: stemcell-2
  vm_type: nkmS20KU9m
  templates:
  - name: foo
    release: foo-release
  - name: bar
    release: bar-release
    notes: active
    quality: superb
`
			yamlMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(yamlString), yamlMap)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return the correct paths", func() {
			paths := pathBuilder.BuildPaths(yamlMap)

			expectedPaths := [][]interface{} {
				{"name"},
				{"releases"},
				{"releases", 0},
				{"releases", 0, "name"},
				{"releases", 0, "version"},
				{"releases", 1},
				{"releases", 1, "name"},
				{"releases", 1, "version"},
				{"stemcells"},
				{"stemcells", 0},
				{"stemcells", 0, "alias"},
				{"stemcells", 0, "name"},
				{"stemcells", 0, "version"},
				{"update"},
				{"update", "canaries"},
				{"update", "canary_watch_time"},
				{"update", "max_in_flight"},
				{"update", "update_watch_time"},
				{"jobs"},
				{"jobs", "name"},
				{"jobs", "instances"},
				{"jobs", "persistent_disk_type"},
				{"jobs", "stemcell"},
				{"jobs", "vm_type"},
				{"jobs", "templates"},
				{"jobs", "templates", 0},
				{"jobs", "templates", 0, "name"},
				{"jobs", "templates", 0, "release"},
				{"jobs", "templates", 1},
				{"jobs", "templates", 1, "name"},
				{"jobs", "templates", 1, "release"},
				{"jobs", "templates", 1, "notes"},
				{"jobs", "templates", 1, "quality"},
			}
			Expect(len(paths)).To(Equal(len(expectedPaths)))
			Expect(paths).To(ConsistOf(expectedPaths))
		})
	})
})
