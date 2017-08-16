package variables_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables/variablesfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Picker", func() {

	var paths [][]interface{}
	var randomizer *variablesfakes.FakeNumberRandomizer
	var pathPicker PathPicker

	BeforeEach(func() {
		randomizer = &variablesfakes.FakeNumberRandomizer{}
		pathPicker = NewPathPicker(randomizer)
	})

	Context("when paths is nil", func() {
		It("returns an empty list of picks", func() {
			picks := pathPicker.PickPaths(nil, 9)
			Expect(picks).ToNot(BeNil())
			Expect(len(picks)).To(Equal(0))
		})
	})

	Context("when paths is empty", func() {
		It("returns an empty list of picks", func() {
			picks := pathPicker.PickPaths([][]interface{}{}, 9)
			Expect(picks).ToNot(BeNil())
			Expect(len(picks)).To(Equal(0))
		})
	})

	Context("when paths contains only sibling entries", func() {
		BeforeEach(func() {
			paths = [][]interface{}{
				{"a", "b", "c"},
				{"d", "e"},
				{"f"},
				{"g", "h"},
				{"i", "j"},
				{"k"},
				{"l"},
				{"m"},
				{"n", "o"},
				{"p", "q", "r"},
				{"s"},
				{"t"},
				{"u"},
				{"v", "w", "x", "y"},
				{"z"},
			}
		})

		Context("when requested picks match the number of provided paths", func() {
			It("returns all of the provided paths", func() {
				picks := pathPicker.PickPaths(paths, len(paths))
				Expect(len(picks)).To(Equal(len(paths)))
				Expect(picks).To(ConsistOf(paths))
			})
		})

		Context("when requested picks exceeds the number of provided paths", func() {
			It("returns all of the provided paths", func() {
				picks := pathPicker.PickPaths(paths, len(paths)*4)
				Expect(len(picks)).To(Equal(len(paths)))
				Expect(picks).To(ConsistOf(paths))
			})
		})

		Context("when requested picks is less than the number of provided paths", func() {
			It("returns the requested number of paths", func() {
				picks := pathPicker.PickPaths(paths, 2)
				Expect(len(picks)).To(Equal(2))
				for _, value := range picks {
					Expect(paths).To(ContainElement(value))
				}
			})

			It("returns all unique paths", func() {
				picks := pathPicker.PickPaths(paths, 2)
				Expect(len(picks)).To(Equal(2))
				for index, value := range picks {
					for otherIndex, otherValue := range picks {
						if otherIndex != index {
							Expect(value).ToNot(Equal(otherValue))
						}
					}
				}
			})

			It("returns randomly selected paths", func() {
				randomizer.IntnReturnsOnCall(0, 8)
				randomizer.IntnReturnsOnCall(1, 4)
				randomizer.IntnReturnsOnCall(2, -1)

				picks := pathPicker.PickPaths(paths, 2)
				Expect(len(picks)).To(Equal(2))
				Expect(picks[0]).To(Equal(paths[8]))
				Expect(picks[1]).To(Equal(paths[4]))
			})
		})
	})

	Context("when paths contains parent and child entries", func() {
		BeforeEach(func() {
			paths = [][]interface{}{
				{"a"},
				{"a", 0},
				{"a", 0, "b"},
				{"a", 0, "c"},
				{"a", 1},
				{"a", 1, "d"},
			}
		})

		Context("when parent is picked", func() {

			BeforeEach(func() {
				randomizer.IntnReturnsOnCall(0, 0)
				randomizer.IntnReturnsOnCall(1, -1)
			})

			It("can not pick its children", func() {
				picks := pathPicker.PickPaths(paths, 2)

				Expect(len(picks)).To(Equal(1))
				Expect(picks[0]).To(Equal(paths[0]))
			})
		})

		Context("when child is picked", func() {

			BeforeEach(func() {
				randomizer.IntnReturnsOnCall(0, 1)
				randomizer.IntnReturnsOnCall(1, 0)
				randomizer.IntnReturnsOnCall(2, -1)
			})

			It("can not pick its ancestors", func() {
				picks := pathPicker.PickPaths(paths, 2)

				Expect(len(picks)).To(Equal(2))
				Expect(picks[0]).To(Equal(paths[1]))
				Expect(picks[1]).To(Equal(paths[4]))
			})
		})
	})
})
