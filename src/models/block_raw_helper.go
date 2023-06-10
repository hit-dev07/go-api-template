package models

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"math"
	"strconv"
)

func ConvertToBlockRaw(value []byte) (*BlockRaw, error) {
	block := BlockRaw{}
	err := protojson.Unmarshal(value, &block)
	if err != nil {
		zap.S().Error("Block_raw_helper: Error in ConvertToBlockRaw: %v", err)
	}
	return &block, err
}

func ValidateHeight(heightString string) bool {
	height, err := strconv.Atoi(heightString)
	if err != nil {
		return false
	}
	if height < 0 || height > math.MaxUint32 {
		return false
	}
	return true
}
