package sdk

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func newFakeClient(l *log.Logger, connected bool, responses [][]byte) *client {
	_ = newClientOptions(&ClientConfig{}, l)
	return &client{
		Client: &fakeMqttClient{
			connected:  connected,
			responses:  responses,
			cmdChan:    make(chan []byte),
			resChan:    make(chan struct{}),
			subscribed: make(map[string][]int),
		},
		logger: l,
	}
}

// fakeMqttClient implements fake mqtt client stub
type fakeMqttClient struct {
	mqtt.Client
	connected bool
	responses [][]byte
	cmdChan   chan []byte
	resChan   chan struct{}
	// published map[string][]int[]byte
	subscribed map[string][]int
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
	var client mqtt.Client
	msg := &fakeMessage{topic: topic}

	c.mockSub(map[string]byte{topic: qos})

	switch toGlobalTopic(topic) {
	case TOPIC_COMMAND:
		go func() {
			select {
			case msg.payload = <-c.cmdChan:
				callback(client, msg)
				c.resChan <- struct{}{}
			case <-time.After(time.Second):
			}
		}()
	case TOPIC_RESPONSE:
		go func() {
			select {
			case <-c.resChan:
				for _, msg.payload = range c.responses {
					time.Sleep(5 * time.Millisecond)
					callback(client, msg)
				}
			case <-time.After(time.Second):
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

		// find the idx inside dictionary
		var idx int
		for i, v := range c.subscribed[gTopic] {
			if v == vin {
				idx = i
				break
			}
		}
		// remove that from dictionary
		c.subscribed[gTopic] = append(c.subscribed[gTopic][:idx], c.subscribed[gTopic][idx+1:]...)
	}

	return &mqtt.DummyToken{}
}

func (c *fakeMqttClient) mockSub(filters map[string]byte) {
	for topic := range filters {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)
		c.subscribed[gTopic] = append(c.subscribed[gTopic], vin)
	}
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

func newFakeResponse(vin int, invoker string, modifier func(*responsePacket)) [][]byte {
	cmd, err := getCmdByInvoker(invoker)
	if err != nil {
		log.Fatal(err)
	}

	// get default rp, and modify it
	rp := newResponsePacket(vin, cmd, nil)
	if modifier != nil {
		modifier(rp)
	}

	// encode
	resBytes, err := encode(rp)
	if err != nil {
		log.Fatal(err)
	}
	if rp.Header.Size == 0 {
		resBytes[2] = uint8(len(resBytes) - 3)
	}

	return [][]byte{
		strToBytes(PREFIX_ACK),
		resBytes,
	}
}

// mockResponse combine response and message to bytes packet.
func mockResponse(r *responsePacket) []byte {
	resBytes, err := encode(&r)
	if err != nil {
		return nil
	}

	// change Header.Size
	if r.Header.Size == 0 {
		resBytes[2] = uint8(len(resBytes) - 3)
	}
	return resBytes
}
