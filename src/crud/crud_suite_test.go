package crud_test

import (
	"github.com/geometry-labs/go-service-template/crud"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	blockRawModel      *crud.BlockRawModel
	blockRawModelMongo *crud.BlockRawModelMongo
)

func TestCrud(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crud Suite")
}

var _ = BeforeSuite(func() {
	blockRawModel = NewBlockModel()
	_ = blockRawModel.Migrate() // Have to create table before running tests
	blockRawModelMongo = NewBlockModelMongo()
})

func NewBlockModel() *crud.BlockRawModel {
	dsn := crud.NewDsn("localhost", "5432", "postgres", "changeme", "test_db", "disable", "UTC")
	postgresConn, _ := crud.NewPostgresConn(dsn)
	testBlockRawModel := crud.NewBlockRawModel(postgresConn.GetConn())
	return testBlockRawModel
}

func NewBlockModelMongo() *crud.BlockRawModelMongo {
	mongoConn := crud.NewMongoConn("mongodb://127.0.0.1:27017")
	blockRawModelMongo := crud.NewBlockRawModelMongo(mongoConn)
	_ = blockRawModelMongo.SetCollectionHandle("icon_test", "contracts")
	return blockRawModelMongo
}
