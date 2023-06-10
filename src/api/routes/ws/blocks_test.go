package ws

import (
	"github.com/geometry-labs/go-service-template/config"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"gopkg.in/Shopify/sarama.v1"

	"github.com/geometry-labs/go-service-template/kafka"
)

func init() {
	//core.GetEnvironment()
	config.Vars.ConfigFile = "config.api.test"
	config.Vars.ConfigType = "yaml"
	config.Vars.ConfigPath = "../../../../envfiles"
	config.ConfigInit()
}

func TestHandlerGetBlocks(t *testing.T) {
	assert := assert.New(t)

	// Create topic broadcaster
	input_chan := make(chan *sarama.ConsumerMessage)
	broadcaster := &kafka.TopicBroadcaster{
		ConsumerChan:   input_chan,
		BroadcastChans: make(map[kafka.BroadcasterID]chan *sarama.ConsumerMessage),
	}
	go broadcaster.Start()

	app := fiber.New()

	app.Use("/", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/", websocket.New(handlerGetBlocks(broadcaster)))
	go app.Listen(":9999")

	test_data := "Test Data"
	go func() {
		for {
			msg := &(sarama.ConsumerMessage{})
			msg.Value = []byte(test_data)

			input_chan <- msg

			time.Sleep(1 * time.Second)
		}
	}()

	// Validate message
	websocket_client, _, err := gorilla.DefaultDialer.Dial("ws://localhost:9999/", nil)
	if err != nil {
		t.Logf("Failed to connect to KafkaWebsocketServer")
		t.Fail()
	}
	defer websocket_client.Close()

	_, message, err := websocket_client.ReadMessage()
	assert.Equal(nil, err)
	assert.Equal(test_data, string(message))

}
