package sdk

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type fakeBroker struct {
	client    mqtt.Client
	responses [][]byte
	cmdChan   chan []byte
	resChan   chan struct{}
}

func (b *fakeBroker) pub(topic string, qos byte, retained bool, payload []byte) error {
	if flush := payload == nil && qos == QOS_CMD_FLUSH; !flush {
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

func (b *fakeBroker) connect() error {
	return nil
}

func (b *fakeBroker) disconnect() {}

func (b *fakeBroker) subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	return nil
}

func (b *fakeBroker) unsubMulti(topics []string) error {
	return nil
}

type fakeMessage struct {
	duplicate bool
	qos       byte
	retained  bool
	topic     string
	messageId uint16
	payload   []byte
}

func (m *fakeMessage) Duplicate() bool {
	return m.duplicate
}
func (m *fakeMessage) Qos() byte {
	return m.qos
}
func (m *fakeMessage) Retained() bool {
	return m.retained
}
func (m *fakeMessage) Topic() string {
	return m.topic
}
func (m *fakeMessage) MessageID() uint16 {
	return m.messageId
}
func (m *fakeMessage) Payload() []byte {
	return m.payload
}
func (m *fakeMessage) Ack() {}
