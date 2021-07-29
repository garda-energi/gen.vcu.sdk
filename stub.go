package sdk

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type packets [][]byte

// stubMqttClient implements stub mqtt client stub
type stubMqttClient struct {
	mqtt.Client
	connected bool

	vinsMutex *sync.RWMutex
	vins      map[int]map[string]packets

	cmdChan map[int]chan []byte
	resChan map[int]chan struct{}
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
	packet := payload.([]byte)
	// feed go routine (command topic) with encoded command
	if flush := packet == nil; !flush {
		vin := getTopicVin(topic)
		c.cmdChan[vin] <- packet
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

		// signal running go routines to stop
		c.vinsMutex.Lock()
		if _, ok := c.vins[vin]; ok {
			switch gTopic {
			case TOPIC_COMMAND:
				close(c.cmdChan[vin])
			case TOPIC_RESPONSE:
				close(c.resChan[vin])
			}

			// remove vin's topic from dictionary
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

		// create topics to vin's dictionary
		c.vinsMutex.Lock()
		if _, ok := c.vins[vin]; !ok {
			c.vins[vin] = make(map[string]packets)
		}
		c.vins[vin][gTopic] = nil
		c.vinsMutex.Unlock()

		// create cmd & response chan (if not initialized)
		switch gTopic {
		case TOPIC_COMMAND, TOPIC_RESPONSE:
			if _, ok := c.cmdChan[vin]; !ok {
				c.cmdChan[vin] = make(chan []byte)
			}
			if _, ok := c.resChan[vin]; !ok {
				c.resChan[vin] = make(chan struct{})
			}
		}

		// make go routines to mock broker behavior on specific topic
		switch gTopic {
		case TOPIC_COMMAND:
			// wait incomming published command, then send signal to (response) go routine
			go func(cmdChan <-chan []byte, resChan chan<- struct{}) {
				for msg := range cmdChan {
					callback(c.Client, &stubMessage{
						topic:   topic,
						payload: msg,
					})
					resChan <- struct{}{}
				}
			}(c.cmdChan[vin], c.resChan[vin])
		case TOPIC_RESPONSE:
			// wait incomming signal from (command) go routine, then pass mock responses to callback
			go func(resChan <-chan struct{}) {
				for range resChan {
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
				}
			}(c.resChan[vin])
		}
	}
}

func (c *stubMqttClient) mockAck(vin int, ack []byte) {
	if ack != nil {
		c.vinsMutex.Lock()
		c.vins[vin][TOPIC_RESPONSE] = packets{ack}
		c.vinsMutex.Unlock()
	}
}

func (c *stubMqttClient) mockResponse(vin int, invoker string, modifier func(*responsePacket)) {
	cmd, err := getCmdByInvoker(invoker)
	if err != nil {
		log.Fatal(err)
	}

	rp := makeResponsePacket(vin, cmd, nil)
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
	c.vins[vin][TOPIC_RESPONSE] = packets{
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
