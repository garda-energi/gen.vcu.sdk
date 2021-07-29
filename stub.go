package sdk

import (
	"fmt"
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

	commandChan  map[int]chan []byte
	responseChan map[int]chan struct{}
	reportChan   map[int]chan [][]byte
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
		c.commandChan[vin] <- packet
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
				close(c.commandChan[vin])
			// case TOPIC_RESPONSE:
			// 	close(c.responseChan[vin])
			case TOPIC_REPORT:
				close(c.reportChan[vin])
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

		// create chan (if not initialized) for specific topic
		switch gTopic {
		case TOPIC_COMMAND, TOPIC_RESPONSE:
			if _, ok := c.commandChan[vin]; !ok {
				c.commandChan[vin] = make(chan []byte)
				c.responseChan[vin] = make(chan struct{})
			}
		case TOPIC_REPORT:
			if _, ok := c.reportChan[vin]; !ok {
				c.reportChan[vin] = make(chan [][]byte)
			}
		}

		// make go routines to mock broker behavior on specific topic
		switch gTopic {
		case TOPIC_COMMAND:
			// wait incomming published command, then send signal to (response) go routine
			go func(commandChan <-chan []byte, responseChan chan<- struct{}) {
				for msg := range commandChan {
					callback(c.Client, &stubMessage{
						topic:   topic,
						payload: msg,
					})
					responseChan <- struct{}{}
				}
				close(responseChan)
			}(c.commandChan[vin], c.responseChan[vin])
		case TOPIC_RESPONSE:
			// wait incomming signal from (command) go routine, then pass mock packets to callback
			go func(responseChan <-chan struct{}) {
				for range responseChan {
					// read packets from vins dictionary
					c.vinsMutex.RLock()
					packets := c.vins[vin][gTopic]
					c.vinsMutex.RUnlock()

					// feed to callback
					for _, msg := range packets {
						time.Sleep(5 * time.Millisecond)
						callback(c.Client, &stubMessage{
							topic:   topic,
							payload: msg,
						})
					}
				}
			}(c.responseChan[vin])
		case TOPIC_REPORT:
			// wait incomming signal, then pass mock packets to callback
			go func(reportChan <-chan [][]byte) {
				for packets := range reportChan {
					fmt.Println(vin, gTopic, packets)

					// feed to callback
					for _, msg := range packets {
						time.Sleep(5 * time.Millisecond)
						callback(c.Client, &stubMessage{
							topic:   topic,
							payload: msg,
						})
					}
				}
			}(c.reportChan[vin])
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

func (c *stubMqttClient) mockReport(vin int, rps []*ReportPacket) {
	// encode fake reports
	res := make(packets, len(rps))
	for i, rp := range rps {
		resBytes, err := encode(rp)
		if err != nil {
			log.Fatal(err)
		}
		res[i] = resBytes
	}

	// tirgger go routine (report) to start
	c.reportChan[vin] <- res
}

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
