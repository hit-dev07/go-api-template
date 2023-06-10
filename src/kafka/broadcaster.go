package kafka

import (
	"gopkg.in/Shopify/sarama.v1"
)

// TODO use uuid for larger ID range
type BroadcasterID int

var LAST_BROADCASTER_ID BroadcasterID = 0

type TopicBroadcastFunc func(channel chan *sarama.ConsumerMessage, message *sarama.ConsumerMessage)

type TopicBroadcaster struct {

	// Input
	ConsumerChan chan *sarama.ConsumerMessage

	// Output
	BroadcastChans map[BroadcasterID]chan *sarama.ConsumerMessage

	BroadcastFunc TopicBroadcastFunc
}

var Broadcasters = map[string]*TopicBroadcaster{}

func newBroadcaster(topic_name string, BroadcastFunc TopicBroadcastFunc) {
	if BroadcastFunc == nil {
		BroadcastFunc = func(channel chan *sarama.ConsumerMessage, message *sarama.ConsumerMessage) {
			channel <- message
		}
	}

	Broadcasters[topic_name] = &TopicBroadcaster{
		make(chan *sarama.ConsumerMessage),
		make(map[BroadcasterID]chan *sarama.ConsumerMessage),
		BroadcastFunc,
	}

	go Broadcasters[topic_name].Start()
}

func (tb *TopicBroadcaster) AddBroadcastChannel(topic_chan chan *sarama.ConsumerMessage) BroadcasterID {
	id := LAST_BROADCASTER_ID
	LAST_BROADCASTER_ID++

	tb.BroadcastChans[id] = topic_chan

	return id
}

func (tb *TopicBroadcaster) RemoveBroadcastChannel(id BroadcasterID) {
	_, ok := tb.BroadcastChans[id]
	if ok {
		delete(tb.BroadcastChans, id)
	}
}

func (tb *TopicBroadcaster) Start() {
	for {
		msg := <-tb.ConsumerChan

		for _, channel := range tb.BroadcastChans {
			tb.BroadcastFunc(channel, msg)
		}
	}
}
