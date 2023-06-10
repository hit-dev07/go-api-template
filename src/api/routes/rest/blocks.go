package rest

import (
	"encoding/json"
	"fmt"
	"github.com/geometry-labs/go-service-template/api/service"
	"github.com/geometry-labs/go-service-template/config"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func BlocksAddHandlers(app *fiber.App) {

	prefix := config.Config.RestPrefix + "/blocks"

	app.Get(prefix+"/", handlerGetQuery)
}

// Blocks
// @Summary Get Blocks that match the query
// @Description Get all blocks in the system.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} *[]models.BlockRaw
// @Router /blocks [get]
func handlerGetQuery(c *fiber.Ctx) error {
	params := new(service.BlocksQueryService)
	if err := c.QueryParser(params); err != nil {
		c.Status(422)
		return c.SendString(fmt.Sprintf("%s", err.Error()))
	}
	if json, err := json.Marshal(params); err == nil {
		zap.S().Debug("params: ", string(json))
	}

	blocks := params.RunQuery(c)
	if len(*blocks) == 0 {
		c.Status(204)
		return c.SendString(fmt.Sprintf("Empty Resultset"))
	}
	zap.S().Debug("ResultSet size:", len(*blocks))

	body, _ := json.Marshal(blocks)
	c.Status(200)
	return c.SendString(string(body))

}
