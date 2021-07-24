package sdk

import (
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type fakeBroker struct {
	client       mqtt.Client
	startPubChan chan struct{}
	command      []byte
	responses    [][]byte
}

func (b *fakeBroker) connect() error {
	return nil
}
func (b *fakeBroker) disconnect() {
}
func (b *fakeBroker) pub(topic string, qos byte, retained bool, payload []byte) error {
	flush := payload == nil && qos == QOS_CMD_FLUSH
	if !flush {
		b.command = payload
		b.startPubChan <- struct{}{}
	}
	return nil
}
func (b *fakeBroker) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	if isSubTopic(TOPIC_RESPONSE, topic) {
		go func() {
			select {
			case <-b.startPubChan:
				min, max := 50, 100
				for _, res := range b.responses {
					rng := rand.Intn(max-min) + min
					time.Sleep(time.Duration(rng) * time.Millisecond)

					handler(b.client, &fakeMessage{
						topic:   topic,
						payload: res,
					})
				}
			case <-time.After(5 * time.Second):
			}
		}()
	}

	return nil
}
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

// type fakeClient struct {
// 	connected bool
// }

// func (c *fakeClient) IsConnected() bool {
// 	return c.connected
// }
// func (c *fakeClient) IsConnectionOpen() bool {
// 	return c.connected
// }
// func (c *fakeClient) Connect() mqtt.Token {
// 	c.connected = true
// 	return &mqtt.DummyToken{}
// }
// func (c *fakeClient) Disconnect(quiesce uint) {
// 	c.connected = false
// }
// func (c *fakeClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
// 	return &mqtt.DummyToken{}
// }
// func (c *fakeClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
// 	return &mqtt.DummyToken{}
// }
// func (c *fakeClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
// 	return &mqtt.DummyToken{}
// }
// func (c *fakeClient) Unsubscribe(topics ...string) mqtt.Token {
// 	return &mqtt.DummyToken{}
// }
// func (c *fakeClient) AddRoute(topic string, callback mqtt.MessageHandler) {

// }
// func (c *fakeClient) OptionsReader() mqtt.ClientOptionsReader {
// 	return mqtt.ClientOptionsReader{}
// }
