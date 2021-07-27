package sdk

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// fakeClient implements fake client stub
type fakeClient struct {
	Client
	connected bool
	responses [][]byte
	cmdChan   chan []byte
	resChan   chan struct{}
	// published map[string][]int[]byte
	subscribed map[string][]int
}

func newFakeClient(connected bool, responses [][]byte) *fakeClient {
	return &fakeClient{
		connected:  connected,
		responses:  responses,
		cmdChan:    make(chan []byte),
		resChan:    make(chan struct{}),
		subscribed: make(map[string][]int),
	}
}

func (c *fakeClient) Connect() mqtt.Token {
	c.connected = true
	return &mqtt.DummyToken{}
}

func (c *fakeClient) Disconnect(quiesce uint) {
	c.connected = false
}

func (c *fakeClient) IsConnected() bool {
	return c.connected
}

func (c *fakeClient) pub(topic string, qos byte, retained bool, payload []byte) error {
	if flush := payload == nil; !flush {
		c.cmdChan <- payload
	}
	return nil
}

func (c *fakeClient) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	var client mqtt.Client
	msg := &fakeMessage{topic: topic}

	c.mockSub(topic)

	switch toGlobalTopic(topic) {
	case TOPIC_COMMAND:
		go func() {
			select {
			case msg.payload = <-c.cmdChan:
				handler(client, msg)
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
					handler(client, msg)
				}
			case <-time.After(time.Second):
			}
		}()
	}
	return nil
}

func (c *fakeClient) subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	c.mockSub(topics...)
	return nil
}

func (c *fakeClient) unsub(topics []string) error {
	c.mockUnsub(topics...)
	return nil
}

func (c *fakeClient) mockSub(topics ...string) {
	for _, topic := range topics {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)
		c.subscribed[gTopic] = append(c.subscribed[gTopic], vin)
	}
}

func (c *fakeClient) mockUnsub(topics ...string) {
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

func fakeResponse(vin int, invoker string) *responsePacket {
	cmd, err := getCmdByInvoker(invoker)
	if err != nil {
		log.Fatal(err)
	}
	return newResponsePacket(vin, cmd, nil)
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
