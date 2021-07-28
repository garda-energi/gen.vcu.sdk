package sdk

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type responses [][]byte

func newStubClient(l *log.Logger, connected bool) *client {
	_ = newClientOptions(&ClientConfig{}, l)
	return &client{
		Client: &stubMqttClient{
			connected: connected,
			cmdChan:   make(chan []byte),
			resChan:   make(chan struct{}),
			stopChan:  make(chan struct{}, 2),
			vins:      make(map[int]map[string]responses),
			vinsMutex: &sync.RWMutex{},
		},
		logger: l,
	}
}

// stubMqttClient implements stub mqtt client stub
type stubMqttClient struct {
	mqtt.Client
	connected bool
	cmdChan   chan []byte
	resChan   chan struct{}
	stopChan  chan struct{}
	vins      map[int]map[string]responses
	vinsMutex *sync.RWMutex
}

func (c *stubMqttClient) Connect() mqtt.Token {
	c.connected = true
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) Disconnect(quiesce uint) {
	c.connected = false
}

func (c *stubMqttClient) IsConnected() bool {
	return c.connected
}

func (c *stubMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	bytes := payload.([]byte)
	if flush := bytes == nil; !flush {
		c.cmdChan <- bytes
	}
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	c.mockSub(map[string]byte{topic: qos}, callback)
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	c.mockSub(filters, callback)
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	for _, topic := range topics {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)

		c.vinsMutex.Lock()
		if _, ok := c.vins[vin]; ok {
			switch gTopic {
			case TOPIC_COMMAND, TOPIC_RESPONSE:
				c.stopChan <- struct{}{}
				c.stopChan <- struct{}{}
			}

			delete(c.vins[vin], gTopic)
		}
		c.vinsMutex.Unlock()
	}

	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) mockSub(filters map[string]byte, callback mqtt.MessageHandler) {
	for topic := range filters {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)

		c.vinsMutex.Lock()
		if _, ok := c.vins[vin]; !ok {
			c.vins[vin] = make(map[string]responses)
		}
		c.vins[vin][gTopic] = nil
		c.vinsMutex.Unlock()

		switch gTopic {
		case TOPIC_COMMAND:
			go func() {
				for {
					select {
					case msg := <-c.cmdChan:
						callback(c.Client, &stubMessage{
							topic:   topic,
							payload: msg,
						})
						c.resChan <- struct{}{}
					case <-c.stopChan:
						return
					}
				}
			}()
		case TOPIC_RESPONSE:
			go func() {
				for {
					select {
					case <-c.resChan:
						c.vinsMutex.RLock()
						responses := c.vins[vin][gTopic]
						c.vinsMutex.RUnlock()

						for _, msg := range responses {
							time.Sleep(5 * time.Millisecond)
							callback(c.Client, &stubMessage{
								topic:   topic,
								payload: msg,
							})
						}
					case <-c.stopChan:
						return
					}
				}
			}()
		}
	}
}

func (c *stubMqttClient) mockAck(vin int, ack []byte) {
	if ack != nil {
		c.vinsMutex.Lock()
		c.vins[vin][TOPIC_RESPONSE] = responses{ack}
		c.vinsMutex.Unlock()
	}
}

func (c *stubMqttClient) mockResponse(vin int, invoker string, modifier func(*responsePacket)) {
	cmd, err := getCmdByInvoker(invoker)
	if err != nil {
		log.Fatal(err)
	}

	rp := newResponsePacket(vin, cmd, nil)
	if modifier != nil {
		modifier(rp)
	}

	resBytes, err := encode(rp)
	if err != nil {
		log.Fatal(err)
	}
	if rp.Header.Size == 0 {
		resBytes[2] = uint8(len(resBytes) - 3)
	}

	c.vinsMutex.Lock()
	c.vins[vin][TOPIC_RESPONSE] = responses{
		strToBytes(PREFIX_ACK),
		resBytes,
	}
	c.vinsMutex.Unlock()
}

// func (c *stubMqttClient) mockReport(vin int, rps []*ReportPacket) {
// 	var res = responses{}
// 	for _, rp := range rps {
// 		resBytes, err := encode(rp)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		res = append(res, resBytes)
// 	}

// 	c.vinsMutex.Lock()
// 	c.vins[vin][TOPIC_REPORT] = res
// 	c.vinsMutex.Unlock()
// }

// stubMessage implements fake message stub
type stubMessage struct {
	mqtt.Message
	topic   string
	payload []byte
}

func (m *stubMessage) Topic() string {
	return m.topic
}
func (m *stubMessage) Payload() []byte {
	return m.payload
}

// stubSleeper implement fake sleeper stub
type stubSleeper struct {
	sleep time.Duration
	after time.Duration
}

func (s *stubSleeper) Sleep(d time.Duration) {
	time.Sleep(s.sleep)
}

func (s *stubSleeper) After(d time.Duration) <-chan time.Time {
	return time.After(s.after)
}
