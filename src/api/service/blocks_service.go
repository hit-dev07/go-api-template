package service

import (
	"github.com/geometry-labs/go-service-template/global"
	"github.com/geometry-labs/go-service-template/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"reflect"
	"strconv"
)

type BlocksQueryService struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`

	Hash      string `query:"hash"`
	Height    uint32 `query:"height"`
	CreatedBy string `query:"created_by"`
	Start     uint32 `query:"start"`
}

func (service *BlocksQueryService) RunQuery(c *fiber.Ctx) *[]models.BlockRaw {
	blocksModel := global.GetGlobal().Blocks
	db := blocksModel.GetDB()

	whereClauseStrings := service.buildWhereClauseStrings()
	orderClauseStrings := service.buildOrderClauseStrings()
	blocks := &[]models.BlockRaw{}
	_ = db.Scopes(Paginate(service)).
		Order(orderClauseStrings).
		Find(blocks, whereClauseStrings...)

	zap.S().Info("Blocks: ", blocks)
	return blocks
}

func (service *BlocksQueryService) buildWhereClauseStrings() []interface{} {
	var strArr []interface{}
	if service.Height > 0 || service.Start > 0 {
		if service.Start > 0 {
			strArr = append(strArr, "number > ?", strconv.Itoa(int(service.Start)))
		} else if service.Height > 0 {
			strArr = append(strArr, "number = ?", strconv.Itoa(int(service.Height)))
		}
	}
	if service.Hash != "" {
		strArr = append(strArr, "hash = ?", service.Hash)
	}
	if service.CreatedBy != "" {
		strArr = append(strArr, "peer_id = ?", service.CreatedBy)
	}
	return strArr
}

func (service *BlocksQueryService) buildOrderClauseStrings() interface{} {
	var strArr string
	strArr = "number desc" //number desc, item_timestamp"
	return strArr
}

func (service *BlocksQueryService) buildLimitClause() int {
	empty := BlocksQueryService{}

	pageSize := 1
	if isEmpty := reflect.DeepEqual(empty, *service); isEmpty {
		return pageSize
	}

	pageSize = service.PageSize
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	return pageSize
}

func Paginate(service *BlocksQueryService) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := service.Page
		if page == 0 {
			page = 1
		}

		pageSize := service.buildLimitClause()

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
