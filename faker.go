package sdk

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type responses [][]byte

func newFakeClient(l *log.Logger, connected bool) *client {
	_ = newClientOptions(&ClientConfig{}, l)
	return &client{
		Client: &fakeMqttClient{
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

// fakeMqttClient implements fake mqtt client stub
type fakeMqttClient struct {
	mqtt.Client
	connected bool
	cmdChan   chan []byte
	resChan   chan struct{}
	stopChan  chan struct{}
	vins      map[int]map[string]responses
	vinsMutex *sync.RWMutex
	// res       responses
	// published map[string][]int[]byte
}

func (c *fakeMqttClient) Connect() mqtt.Token {
	c.connected = true
	return &mqtt.DummyToken{}
}

func (c *fakeMqttClient) Disconnect(quiesce uint) {
	c.connected = false
}

func (c *fakeMqttClient) IsConnected() bool {
	return c.connected
}

func (c *fakeMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	bytes := payload.([]byte)
	if flush := bytes == nil; !flush {
		c.cmdChan <- bytes
	}
	return &mqtt.DummyToken{}
}

func (c *fakeMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	c.mockSub(map[string]byte{topic: qos})

	gTopic := toGlobalTopic(topic)
	vin := getTopicVin(topic)
	msg := &fakeMessage{topic: topic}

	switch toGlobalTopic(topic) {
	case TOPIC_COMMAND:
		go func() {
			for {
				select {
				case msg.payload = <-c.cmdChan:
					callback(c.Client, msg)
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

					for _, msg.payload = range responses {
						time.Sleep(5 * time.Millisecond)
						callback(c.Client, msg)
					}
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	return &mqtt.DummyToken{}
}

func (c *fakeMqttClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	c.mockSub(filters)
	return &mqtt.DummyToken{}
}

func (c *fakeMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	for _, topic := range topics {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)

		c.vinsMutex.RLock()
		_, ok := c.vins[vin]
		c.vinsMutex.RUnlock()

		if ok {
			switch gTopic {
			case TOPIC_COMMAND, TOPIC_RESPONSE:
				c.stopChan <- struct{}{}
				c.stopChan <- struct{}{}
			}

			c.vinsMutex.Lock()
			delete(c.vins[vin], gTopic)
			c.vinsMutex.Unlock()
		}
	}

	return &mqtt.DummyToken{}
}

func (c *fakeMqttClient) mockSub(filters map[string]byte) {
	for topic := range filters {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)

		c.vinsMutex.RLock()
		_, ok := c.vins[vin]
		c.vinsMutex.RUnlock()

		c.vinsMutex.Lock()
		if !ok {
			c.vins[vin] = make(map[string]responses)
		}
		c.vins[vin][gTopic] = nil
		c.vinsMutex.Unlock()
	}
}

func (c *fakeMqttClient) mockAck(vin int, ack []byte) {
	if ack == nil {
		return
	}
	c.vinsMutex.Lock()
	c.vins[vin][TOPIC_RESPONSE] = responses{ack}
	c.vinsMutex.Unlock()
}

func (c *fakeMqttClient) mockResponse(vin int, invoker string, modifier func(*responsePacket)) {
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

func (c *fakeMqttClient) mockReport(vin int, rps []*ReportPacket) {
	var res = responses{}
	for _, rp := range rps {
		resBytes, err := encode(rp)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, resBytes)
	}

	c.vinsMutex.Lock()
	c.vins[vin][TOPIC_REPORT] = res
	c.vinsMutex.Unlock()
}

// fakeMessage implements fake message stub
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

// fakeSleeper implement fake sleeper stub
type fakeSleeper struct {
	sleep time.Duration
	after time.Duration
}

func (s *fakeSleeper) Sleep(d time.Duration) {
	time.Sleep(s.sleep)
}

func (s *fakeSleeper) After(d time.Duration) <-chan time.Time {
	return time.After(s.after)
}
