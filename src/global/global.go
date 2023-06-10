package global

import (
	"sync"

	"github.com/geometry-labs/go-service-template/crud"
)

const Version = "v0.1.0"

type Global struct {
	Blocks *crud.BlockRawModel
}

var globalInstance *Global
var globalOnce sync.Once

func GetGlobal() *Global {
	globalOnce.Do(func() {
		globalInstance = &Global{
			Blocks: crud.GetBlockRawModel(),
		}
	})
	return globalInstance
}

var ShutdownChan = make(chan int)
