package sdk

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type fakeBroker struct {
	Broker
	client    mqtt.Client
	responses [][]byte
	cmdChan   chan []byte
	resChan   chan struct{}
}

func (b *fakeBroker) pub(topic string, qos byte, retained bool, payload []byte) error {
	if flush := payload == nil; !flush {
		b.cmdChan <- payload
	}
	return nil
}

func (b *fakeBroker) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	switch toGlobalTopic(topic) {
	case TOPIC_COMMAND:
		go func() {
			select {
			case cmdPacket := <-b.cmdChan:
				handler(b.client, &fakeMessage{
					topic:   topic,
					payload: cmdPacket,
				})
				b.resChan <- struct{}{}
			case <-time.After(time.Second):
			}
		}()
	case TOPIC_RESPONSE:
		go func() {
			select {
			case <-b.resChan:
				for _, resPacket := range b.responses {
					randomSleep(50, 100, time.Millisecond)
					handler(b.client, &fakeMessage{
						topic:   topic,
						payload: resPacket,
					})
				}
			case <-time.After(time.Second):
			}
		}()
	}
	return nil
}

func (b *fakeBroker) unsubMulti(topics []string) error {
	return nil
}

type fakeMessage struct {
	mqtt.Message
	topic   string
	payload []byte
}

func (m *fakeMessage) Topic() string {
	return m.topic
}
func (m *fakeMessage) Payload() []byte {
	return m.payload
}
