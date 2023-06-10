package loader

import (
	"fmt"
	"github.com/geometry-labs/go-service-template/global"
	"github.com/geometry-labs/go-service-template/models"
	"go.uber.org/zap"
)

func StartBlockRawsLoader() {
	go BlockRawsLoader()
}

func BlockRawsLoader() {
	var block *models.BlockRaw
	postgresLoaderChan := global.GetGlobal().Blocks.GetWriteChan()
	for {
		block = <-postgresLoaderChan
		global.GetGlobal().Blocks.RetryCreate(block) // inserted here !!
		zap.S().Debug(fmt.Sprintf(
			"Loader BlockRaws: Loaded in postgres table BlockRaws, Block Number %d", block.Number),
		)
	}
}
