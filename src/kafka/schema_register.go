package kafka

import (
	"encoding/binary"
	"github.com/geometry-labs/go-service-template/config"
	"go.uber.org/zap"
	"io/ioutil"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/riferrei/srclient"
)

type RegisterSchemaFunc func(topic string, isKey bool, srcSchemaFile string, forceUpdate bool) (int, error)

func RegisterSchema(topic string, isKey bool, srcSchemaFile string, forceUpdate bool) (int, error) {
	zap.S().Info("RegisterSchema() \n")
	schemaRegistryClient := srclient.CreateSchemaRegistryClient("http://" + config.Config.SchemaRegistryURL)
	schema, err := schemaRegistryClient.GetLatestSchema(topic, false)
	if schema == nil {
		schema, err = registerSchema(schemaRegistryClient, topic, isKey, srcSchemaFile)
	} else if forceUpdate {
		schema, err = registerSchema(schemaRegistryClient, topic, isKey, srcSchemaFile) //TODO: Resolve update not happening
	}

	if err != nil {
		return 0, err
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	return schema.ID(), nil
}

func registerSchema(schemaRegistryClient *srclient.SchemaRegistryClient, topic string, isKey bool, srcSchemaFile string) (*srclient.Schema, error) {
	filePath := "schemas/" + srcSchemaFile + ".proto"
	zap.S().Info("Adding/Updating Schema from: %s\n", filePath)
	schemaBytes, _ := ioutil.ReadFile(filePath)
	schema, err := schemaRegistryClient.CreateSchema(topic, string(schemaBytes), srclient.Protobuf, isKey)
	if err != nil {
		//panic(fmt.Sprintf("Error creating the schema %s", err))
		zap.S().Info("Error creating the schema %s\n", err)
		return nil, err
	}
	return schema, nil
}

func RetriableRegisterSchema(fn RegisterSchemaFunc, topic string, isKey bool, srcSchemaFile string, forceUpdate bool) (int, error) {
	x := 0
	operation := func() error {
		val, err := fn(topic, isKey, srcSchemaFile, forceUpdate)
		if err != nil {
			zap.S().Info("RegisterSchema unsuccessful")
		} else {
			x = val
		}
		return err
	}
	neb := backoff.NewExponentialBackOff()
	neb.MaxElapsedTime = time.Minute
	err := backoff.Retry(operation, neb)
	if err != nil {
		zap.S().Info("Finally also RegisterSchema Unsuccessful")
	} else {
		zap.S().Info("Finally RegisterSchema Successful")
	}
	return x, err
}
