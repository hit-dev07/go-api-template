package ws

import (
	"github.com/geometry-labs/go-service-template/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gopkg.in/Shopify/sarama.v1"

	"github.com/geometry-labs/go-service-template/kafka"
)

func BlocksAddHandlers(app *fiber.App) {

	prefix := config.Config.WebsocketPrefix + "/blocks"

	app.Use(prefix, func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get(prefix+"/", websocket.New(handlerGetBlocks(kafka.Broadcasters["blocks"])))
}

func handlerGetBlocks(broadcaster *kafka.TopicBroadcaster) func(c *websocket.Conn) {

	return func(c *websocket.Conn) {

		// Add broadcaster
		topic_chan := make(chan *sarama.ConsumerMessage)
		id := broadcaster.AddBroadcastChannel(topic_chan)
		defer func() {
			// Remove broadcaster
			broadcaster.RemoveBroadcastChannel(id)
		}()

		// Read for close
		client_close_sig := make(chan bool)
		go func() {
			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					client_close_sig <- true
					break
				}
			}
		}()

		for {
			// Read
			msg := <-topic_chan

			// Broadcast
			err := c.WriteMessage(websocket.TextMessage, msg.Value)
			if err != nil {
				break
			}

			// check for client close
			select {
			case _ = <-client_close_sig:
				break
			default:
				continue
			}
		}
	}
}
