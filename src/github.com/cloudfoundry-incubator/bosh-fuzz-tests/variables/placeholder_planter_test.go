package variables_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("PlaceholderPlanter", func() {

	var manifestString string
	var manifestMap map[interface{}]interface{}
	var candidates [][]interface{}
	var nameGenerator fakes.FakeNameGenerator
	var placeholderPlanter PlaceholderPlanter

	Context("when the manifest contains simple elements", func() {

		BeforeEach(func() {
			manifestString = `
version: nothing
release: fake
name: foo-deployment
`
			manifestMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(manifestString), manifestMap)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the candidates are valid paths", func() {
			BeforeEach(func() {
				candidates = [][]interface{}{{"version"}, {"release"}}
				nameGenerator = fakes.FakeNameGenerator{
					Names: []string{"placeholder1", "placeholder2"},
				}
				placeholderPlanter = NewPlaceholderPlanter(&nameGenerator)
			})

			It("should update the manifest", func() {
				_, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())
				Expect(manifestMap["version"]).To(Equal("((placeholder1))"))
				Expect(manifestMap["release"]).To(Equal("((placeholder2))"))
			})

			It("should return the placeholder map with the substituted values", func() {
				substitutions, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())
				Expect(substitutions["placeholder1"]).To(Equal("nothing"))
				Expect(substitutions["placeholder2"]).To(Equal("fake"))
			})
		})
	})

	Context("when the manifest contains child elements", func() {

		BeforeEach(func() {
			manifestString = `
version:
  release: fake
  name: foo-deployment
`
			manifestMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(manifestString), manifestMap)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the candidates are valid paths", func() {
			BeforeEach(func() {
				candidates = [][]interface{}{{"version", "release"}, {"version", "name"}}
				nameGenerator = fakes.FakeNameGenerator{
					Names: []string{"placeholder1", "placeholder2"},
				}
				placeholderPlanter = NewPlaceholderPlanter(&nameGenerator)
			})

			It("should update the manifest", func() {
				_, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())

				versionMap := manifestMap["version"].(map[interface{}]interface{})
				Expect(versionMap["release"]).To(Equal("((placeholder1))"))
				Expect(versionMap["name"]).To(Equal("((placeholder2))"))
			})

			It("should return the placeholder map with the substituted values", func() {
				substitutions, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())
				Expect(substitutions["placeholder1"]).To(Equal("fake"))
				Expect(substitutions["placeholder2"]).To(Equal("foo-deployment"))
			})
		})
	})

	Context("when the manifest contains array elements", func() {

		BeforeEach(func() {
			manifestString = `
version: [1, 2, 3]
`
			manifestMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(manifestString), manifestMap)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the candidates are valid paths", func() {
			BeforeEach(func() {
				candidates = [][]interface{}{{"version", 0}, {"version", 2}}
				nameGenerator = fakes.FakeNameGenerator{
					Names: []string{"placeholder1", "placeholder2"},
				}
				placeholderPlanter = NewPlaceholderPlanter(&nameGenerator)
			})

			It("should update the manifest", func() {
				_, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())

				versionArray := manifestMap["version"].([]interface{})
				Expect(versionArray[0]).To(Equal("((placeholder1))"))
				Expect(versionArray[1]).To(Equal(2))
				Expect(versionArray[2]).To(Equal("((placeholder2))"))
			})

			It("should return the placeholder map with the substituted values", func() {
				substitutions, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())
				Expect(substitutions["placeholder1"]).To(Equal(1))
				Expect(substitutions["placeholder2"]).To(Equal(3))
			})
		})
	})

	Context("when the manifest contains simple, nested and array elements", func() {

		BeforeEach(func() {
			manifestString = `
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
			manifestMap = map[interface{}]interface{}{}
			err := yaml.Unmarshal([]byte(manifestString), manifestMap)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the candidates are valid paths", func() {
			BeforeEach(func() {
				candidates = [][]interface{}{
					{"name"},
					{"update"},
					{"releases", 1, "version"},
					{"stemcells", 0},
					{"jobs", "templates"},
				}
				nameGenerator = fakes.FakeNameGenerator{
					Names: []string{"placeholder1", "placeholder2", "placeholder3", "placeholder4", "placeholder5"},
				}
				placeholderPlanter = NewPlaceholderPlanter(&nameGenerator)
			})

			It("should update the manifest", func() {
				_, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())

				Expect(manifestMap["name"].(string)).To(Equal("((placeholder1))"))
				Expect(manifestMap["update"].(string)).To(Equal("((placeholder2))"))

				releasesArray := manifestMap["releases"].([]interface{})
				elementMap := releasesArray[1].(map[interface{}]interface{})
				Expect(elementMap["version"]).To(Equal("((placeholder3))"))

				stemcellsArray := manifestMap["stemcells"].([]interface{})
				Expect(stemcellsArray[0]).To(Equal("((placeholder4))"))

				instanceGroupsMap := manifestMap["jobs"].(map[interface{}]interface{})
				Expect(instanceGroupsMap["templates"].(string)).To(Equal("((placeholder5))"))
			})

			It("should return the placeholder map with the substituted values", func() {
				substitutions, err := placeholderPlanter.PlantPlaceholders(&manifestMap, candidates)
				Expect(err).ToNot(HaveOccurred())

				// {"name"}
				Expect(substitutions["placeholder1"]).To(Equal("foo-deployment"))

				// {"update"}
				valueMap := map[interface{}]interface{}{}
				valueMap["canaries"] = 2
				valueMap["canary_watch_time"] = 4000
				valueMap["max_in_flight"] = 2
				valueMap["update_watch_time"] = 20
				Expect(substitutions["placeholder2"]).To(Equal(valueMap))

				// {"releases", 1, "version"},
				Expect(substitutions["placeholder3"]).To(Equal("30+dev.6"))

				// {"stemcells", 0},
				valueMap = map[interface{}]interface{}{}
				valueMap["alias"] = "stemcell-2"
				valueMap["name"] = "ubuntu-stemcell"
				valueMap["version"] = "latest"
				Expect(substitutions["placeholder4"]).To(Equal(valueMap))

				// {"jobs", "templates"}
				valueMap1 := map[interface{}]interface{}{}
				valueMap1["name"] = "foo"
				valueMap1["release"] = "foo-release"

				valueMap2 := map[interface{}]interface{}{}
				valueMap2["name"] = "bar"
				valueMap2["release"] = "bar-release"
				valueMap2["notes"] = "active"
				valueMap2["quality"] = "superb"

				valueArray := []interface{}{valueMap1, valueMap2}
				Expect(substitutions["placeholder5"]).To(Equal(valueArray))
			})
		})
	})
})
