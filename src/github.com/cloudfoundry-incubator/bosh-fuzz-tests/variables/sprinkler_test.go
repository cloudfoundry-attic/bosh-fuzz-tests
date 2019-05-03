package variables_test

import (
	"errors"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables/variablesfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Sprinkler", func() {

	Context("when called with a manifest", func() {

		var parameters bftconfig.Parameters
		var fs *fakesys.FakeFileSystem
		var sprinkler Sprinkler
		var randomizer *FakeNumberRandomizer
		var pathBuilder *FakePathBuilder
		var pathPicker *FakePathPicker
		var placeholderPlanter *FakePlaceholderPlanter
		var nameGenerator name_generator.NameGenerator

		BeforeEach(func() {
			parameters = bftconfig.Parameters{NumOfSubstitutions: []int{7}}
			fs = fakesys.NewFakeFileSystem()
			randomizer = &FakeNumberRandomizer{}
			pathBuilder = &FakePathBuilder{}
			pathPicker = &FakePathPicker{}
			placeholderPlanter = &FakePlaceholderPlanter{}
			nameGenerator = &fakes.FakeNameGenerator{
				Names: []string{"placeholder1", "placeholder2"},
			}

			sprinkler = NewSprinkler(parameters, fs, randomizer, pathBuilder, pathPicker, placeholderPlanter, nameGenerator)
		})

		Context("when manifest exists at given path", func() {
			BeforeEach(func() {
				fs.WriteFile("manifest-path", []byte("---\n"))
			})

			It("does NOT raise an error", func() {
				_, err := sprinkler.SprinklePlaceholders("manifest-path")
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when manifest is invalid", func() {
				var yamlString string
				BeforeEach(func() {
					yamlString = `bad-content`
					fs.WriteFile("manifest-path", []byte(yamlString))
				})

				It("raises an error", func() {
					_, err := sprinkler.SprinklePlaceholders("manifest-path")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Error unmarshalling manifest file"))
				})
			})

			Context("when manifest is valid", func() {
				var yamlString string

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

instance_groups:
  name: zRD
  instances: 5
  persistent_disk_type: czcuBXB7WY
  stemcell: stemcell-2
  vm_type: nkmS20KU9m
  jobs:
  - name: foo
    release: foo-release
  - name: bar
    release: bar-release
    notes: active
    quality: superb
`
					bytes := []byte(yamlString)
					fs.WriteFile("manifest-path", bytes)
				})

				It("updates the manifest with placeholder", func() {
					placeholderPlanter.PlantPlaceholdersStub = func(manifest *map[interface{}]interface{}, candidates [][]interface{}) (map[string]interface{}, error) {
						(*manifest)["name"] = "((placeholder1))"
						return nil, nil
					}

					_, err := sprinkler.SprinklePlaceholders("manifest-path")
					Expect(err).ToNot(HaveOccurred())

					updatedYamlString, _ := fs.ReadFile("manifest-path")
					updatedManifest := map[interface{}]interface{}{}
					yaml.Unmarshal([]byte(updatedYamlString), updatedManifest)
					Expect(updatedManifest["name"]).To(Equal("((placeholder1))"))
				})

				It("returns the placeholder map", func() {
					placeholderPlanter.PlantPlaceholdersStub = func(manifest *map[interface{}]interface{}, candidates [][]interface{}) (map[string]interface{}, error) {
						return map[string]interface{}{"placeholder1": "foo-deployment"}, nil
					}

					result, err := sprinkler.SprinklePlaceholders("manifest-path")
					Expect(err).ToNot(HaveOccurred())

					Expect(result).To(Equal(map[string]interface{}{"placeholder1": "foo-deployment"}))
				})

				Context("when planting placeholders errors", func() {
					BeforeEach(func() {
						placeholderPlanter.PlantPlaceholdersReturns(nil, errors.New("sample-error"))
					})

					It("raises an error", func() {
						_, err := sprinkler.SprinklePlaceholders("manifest-path")
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Error adding variables to manifest file: sample-error"))
					})
				})

				Context("when writing manifest with placeholder to file fails", func() {
					BeforeEach(func() {
						fs.WriteFileError = errors.New("write-error")
					})

					It("raises an error", func() {
						_, err := sprinkler.SprinklePlaceholders("manifest-path")
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Error writing manifest file: write-error"))
					})
				})
			})
		})

		Context("when manifest does NOT exist at given path", func() {
			BeforeEach(func() {
				fs.ReadFileError = errors.New("error")
			})

			It("raises an error", func() {
				_, err := sprinkler.SprinklePlaceholders("bad-manifest-path")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Error reading manifest file: Not found: open bad-manifest-path: no such file or directory"))
			})
		})
	})
})
