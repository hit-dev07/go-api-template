package fixtures

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"

	"github.com/geometry-labs/go-service-template/models"
)

const (
	Block_raws_fixture = "block_raws.json"
)

type Fixtures []Fixture
type Fixture struct {
	Input    map[string]interface{}
	Expected map[string]interface{}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadTestFixtures(file string) (Fixtures, error) {
	var fs Fixtures
	dat, err := ioutil.ReadFile(getFixtureDir() + file)
	check(err)
	err = json.Unmarshal(dat, &fs)
	return fs, err
}

func getFixtureDir() string {
	callDir, _ := os.Getwd()
	callDirSplit := strings.Split(callDir, "/")
	for i := len(callDirSplit) - 1; i >= 0; i-- {
		if callDirSplit[i] != "src" {
			callDirSplit = callDirSplit[:len(callDirSplit)-1]
		} else {
			break
		}
	}
	callDirSplit = append(callDirSplit, "fixtures")
	fixtureDir := strings.Join(callDirSplit, "/")
	fixtureDir = fixtureDir + "/"
	zap.S().Info(fixtureDir)
	return fixtureDir
}

func ReadCurrentDir() {
	file, err := os.Open(".")
	if err != nil {
		zap.S().Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		zap.S().Info(name)
	}
}

func (f *Fixture) GetBlock(data map[string]interface{}) *models.BlockRaw {
	block := models.BlockRaw{
		Signature:        data["signature"].(string),
		ItemId:           data["item_id"].(string),
		NextLeader:       data["next_leader"].(string),
		TransactionCount: uint32(data["transaction_count"].(float64)),
		Type:             data["type"].(string),
		Version:          data["version"].(string),
		PeerId:           data["peer_id"].(string),
		Number:           uint32(data["number"].(float64)),
		MerkleRootHash:   data["merkle_root_hash"].(string),
		ItemTimestamp:    data["item_timestamp"].(string),
		Hash:             data["hash"].(string),
		ParentHash:       data["parent_hash"].(string),
		Timestamp:        uint64(data["timestamp"].(float64)),
	}
	return &block
}
