package crud

import (
	"github.com/geometry-labs/go-service-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"sync"
)

type BlockRawModelMongo struct {
	mongoConn *MongoConn
	model     *models.BlockRaw
	//databaseHandle   *mongo.Database
	collectionHandle *mongo.Collection
}

type KeyValue struct {
	Key   string
	Value interface{}
}

var blockRawModelMongoInstance *BlockRawModelMongo
var blockRawModelMongoOnce sync.Once

func GetBlockRawModelMongo() *BlockRawModelMongo {
	blockRawModelMongoOnce.Do(func() {
		blockRawModelMongoInstance = &BlockRawModelMongo{
			mongoConn: GetMongoConn(),
			model:     &models.BlockRaw{},
		}
	})
	return blockRawModelMongoInstance
}

func NewBlockRawModelMongo(conn *MongoConn) *BlockRawModelMongo {
	blockRawModelMongoInstance := &BlockRawModelMongo{
		mongoConn: conn,
		model:     &models.BlockRaw{},
	}
	return blockRawModelMongoInstance
}

//func (b *BlockRawModelMongo) SetCollectionHandle(collection *mongo.Collection) {
//	b.collectionHandle = collection
//}

func (b *BlockRawModelMongo) GetMongoConn() *MongoConn {
	return b.mongoConn
}

func (b *BlockRawModelMongo) GetModel() *models.BlockRaw {
	return b.model
}

func (b *BlockRawModelMongo) SetCollectionHandle(database string, collection string) *mongo.Collection {
	b.collectionHandle = b.mongoConn.DatabaseHandle(database).Collection(collection)
	return b.collectionHandle
}

func (b *BlockRawModelMongo) GetCollectionHandle() *mongo.Collection {
	return b.collectionHandle
}

func (b *BlockRawModelMongo) InsertOne(block *models.BlockRaw) (*mongo.InsertOneResult, error) {
	one, err := b.collectionHandle.InsertOne(b.mongoConn.ctx, block)
	return one, err
}

func (b *BlockRawModelMongo) DeleteMany(kv *KeyValue) (*mongo.DeleteResult, error) {
	delR, err := b.collectionHandle.DeleteMany(b.mongoConn.ctx, bson.D{{kv.Key, kv.Value}})
	return delR, err
}

func (b *BlockRawModelMongo) find(kv *KeyValue) (*mongo.Cursor, error) {
	cursor, err := b.collectionHandle.Find(b.mongoConn.ctx, bson.D{{kv.Key, kv.Value}})
	return cursor, err
}

func (b *BlockRawModelMongo) FindAll(kv *KeyValue) []bson.M {
	cursor, err := b.find(kv)
	if err != nil {
		zap.S().Info("Exception in getting a curser to a find in mongodb: ", err)
	}
	var results []bson.M
	if err = cursor.All(b.mongoConn.ctx, &results); err != nil {
		zap.S().Info("Exception in find all: ", err)
	}
	return results

}
