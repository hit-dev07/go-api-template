//+build integration

package crud_test

import (
	"github.com/geometry-labs/go-service-template/crud"
	"github.com/geometry-labs/go-service-template/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mongo Integration test", func() {
	testFixtures, _ := fixtures.LoadTestFixtures(fixtures.Block_raws_fixture)

	Describe("BlockModel with mongodb", func() {

		Context("Insert in block collection", func() {
			//testFixtures, _ = fixtures.LoadTestFixtures(fixtures.Block_raws_fixture) //To
			for _, fixture := range testFixtures {
				block := fixture.GetBlock(fixture.Input)
				kv := &crud.KeyValue{
					Key:   "signature",
					Value: block.Signature,
				}
				BeforeEach(func() {
					blockRawModelMongo.DeleteMany(&crud.KeyValue{Key: "signature", Value: block.Signature})
				})
				It("insert in mongodb", func() {
					_, err := blockRawModelMongo.InsertOne(block)
					if err != nil {
						Expect(1).To(Equal(0))
					}
					Expect(err).To(BeNil())

					// test find
					results := blockRawModelMongo.FindAll(kv)
					Expect(len(results) == 1).To(Equal(true))

					//test delete
					_, err = blockRawModelMongo.DeleteMany(kv)
					Expect(err).To(BeNil())
				}) // It
			} // For each fixture
		}) // Context "Insert in block collection"

	}) // Describe "BlockModel with mongodb"

}) // Describe "BlockModelMongo"
