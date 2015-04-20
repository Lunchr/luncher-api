package db_test

import (
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
)

var _ = Describe("Tags", func() {
	Describe("GetName", func() {
		It("should get tag by name", func(done Done) {
			defer close(done)
			tag, err := tagsCollection.GetName("kala")
			Expect(err).NotTo(HaveOccurred())
			Expect(tag.Name).To(Equal("kala"))
			Expect(tag.DisplayName).To(Equal("Kalast"))
		})

		It("should get tag by name", func(done Done) {
			defer close(done)
			tag, err := tagsCollection.GetName("siga")
			Expect(err).NotTo(HaveOccurred())
			Expect(tag.Name).To(Equal("siga"))
			Expect(tag.DisplayName).To(Equal("Seast"))
		})

		It("should return nothing if doesn't exist", func(done Done) {
			defer close(done)
			_, err := tagsCollection.GetName("blablabla")
			Expect(err).To(Equal(mgo.ErrNotFound))
		})
	})

	Describe("GetAll", func() {
		It("should list all the tags", func(done Done) {
			defer close(done)
			iter := tagsCollection.GetAll()
			count := 0
			tagNames := map[string]int{
				"kala":   0,
				"lind":   0,
				"siga":   0,
				"veis":   0,
				"lammas": 0,
			}
			var tag model.Tag
			for iter.Next(&tag) {
				Expect(tag).NotTo(BeNil())
				i, contains := tagNames[tag.Name]
				Expect(contains).To(BeTrue())
				Expect(i).To(Equal(0))
				tagNames[tag.Name]++
				count++
			}
			Expect(count).To(Equal(5))
			err := iter.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("UpdateName", func() {
		RebuildDBAfterEach()
		It("should fail for a non-existent name", func(done Done) {
			defer close(done)
			err := tagsCollection.UpdateName("a random name", &model.Tag{})
			Expect(err).To(HaveOccurred())
		})

		Context("with atag with known ID inserted", func() {
			var name string
			BeforeEach(func(done Done) {
				defer close(done)
				name = "a test name"
				err := tagsCollection.Insert(&model.Tag{
					Name: name,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should update the tag in DB", func(done Done) {
				defer close(done)
				err := tagsCollection.UpdateName(name, &model.Tag{
					Name: "an updated name",
				})
				Expect(err).NotTo(HaveOccurred())
				tag, err := tagsCollection.GetName("an updated name")
				Expect(err).NotTo(HaveOccurred())
				Expect(tag.Name).To(Equal("an updated name"))
			})
		})
	})
})
