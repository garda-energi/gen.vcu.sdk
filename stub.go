package sdk

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type resChan chan struct{}
type cmdChan chan packet
type stsChan chan packet
type repChan chan packets

// stubMqttClient implements stub mqtt client stub
type stubMqttClient struct {
	mqtt.Client
	connected bool

	responses *sync.Map

	ch struct {
		res *sync.Map
		cmd *sync.Map
		rep *sync.Map
		sts *sync.Map
	}
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
	packet := payload.(packet)
	if flush := packet == nil; !flush {
		vin := getTopicVin(topic)
		// feed go routine (command topic) with encoded command
		if ch, ok := c.ch.cmd.Load(vin); ok {
			ch.(cmdChan) <- packet
		}
	}
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	c.stubSub(map[string]byte{topic: qos}, callback)
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	c.stubSub(filters, callback)
	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) stubSub(filters map[string]byte, callback mqtt.MessageHandler) {
	for topic := range filters {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)

		// make go routines to mock broker behavior on specific topic
		switch gTopic {
		case TOPIC_COMMAND:
			if _, ok := c.ch.cmd.Load(vin); !ok {
				c.ch.cmd.Store(vin, make(cmdChan))
				c.ch.res.Store(vin, make(resChan))
			}

			// wait incomming published command, then send signal to (response) go routine
			go func(vin int, topic string) {
				chCmd, _ := c.ch.cmd.Load(vin)
				chRes, _ := c.ch.res.Load(vin)

				for msg := range chCmd.(cmdChan) {
					callback(c.Client, &stubMessage{
						topic:   topic,
						payload: msg,
					})

					chRes.(resChan) <- struct{}{}
				}
				// end command
				c.ch.cmd.Delete(vin)
				close(chRes.(resChan))
			}(vin, topic)
		case TOPIC_RESPONSE:
			c.responses.Store(vin, packets{})

			// wait incomming signal from (command) go routine, then pass mock packets to callback
			go func(vin int, topic string) {
				chRes, _ := c.ch.res.Load(vin)

				for range chRes.(resChan) {
					// read packets from vins dictionary
					if res, ok := c.responses.Load(vin); ok {
						// feed to callback
						for _, msg := range res.(packets) {
							time.Sleep(5 * time.Millisecond)
							callback(c.Client, &stubMessage{
								topic:   topic,
								payload: msg,
							})
						}
					}
				}
				// end response
				c.ch.res.Delete(vin)
				c.responses.Delete(vin)
			}(vin, topic)
		case TOPIC_REPORT:
			if _, ok := c.ch.rep.Load(vin); !ok {
				c.ch.rep.Store(vin, make(repChan))
			}

			// wait incomming signal, then pass mock packets to callback
			go func(vin int, topic string) {
				chRep, _ := c.ch.rep.Load(vin)

				for packets := range chRep.(repChan) {
					// feed to callback
					for _, msg := range packets {
						time.Sleep(5 * time.Millisecond)
						callback(c.Client, &stubMessage{
							topic:   topic,
							payload: msg,
						})
					}
				}
				// end of command
				c.ch.rep.Delete(vin)
			}(vin, topic)
		case TOPIC_STATUS:
			if _, ok := c.ch.sts.Load(vin); !ok {
				c.ch.sts.Store(vin, make(stsChan))
			}

			// wait incomming signal, then pass mock packets to callback
			go func(vin int, topic string) {
				chSts, _ := c.ch.sts.Load(vin)

				for msg := range chSts.(stsChan) {
					callback(c.Client, &stubMessage{
						topic:   topic,
						payload: msg,
					})
				}
				// end of status
				c.ch.sts.Delete(vin)
			}(vin, topic)
		}
	}
}

func (c *stubMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	for _, topic := range topics {
		gTopic := toGlobalTopic(topic)
		vin := getTopicVin(topic)

		// signal running go routines to stop
		switch gTopic {
		case TOPIC_COMMAND:
			if ch, ok := c.ch.cmd.Load(vin); ok {
				close(ch.(cmdChan))
			}
		case TOPIC_REPORT:
			if ch, ok := c.ch.rep.Load(vin); ok {
				close(ch.(repChan))
			}
		case TOPIC_STATUS:
			if ch, ok := c.ch.sts.Load(vin); ok {
				close(ch.(stsChan))
			}
		}
	}

	return &mqtt.DummyToken{}
}

func (c *stubMqttClient) mockAck(vin int, ack packet) {
	if ack != nil {
		c.responses.Store(vin, packets{ack})
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

	resBytes, err := encodePacket(rp)
	if err != nil {
		log.Fatal(err)
	}

	c.responses.Store(vin, packets{
		strToBytes(PREFIX_ACK),
		resBytes,
	})
}

func (c *stubMqttClient) mockReports(vin int, rps []*ReportPacket) {
	// encode fake reports
	res := make(packets, len(rps))
	for i, rp := range rps {
		resBytes, err := encodePacket(rp)
		if err != nil {
			log.Fatal(err)
		}
		res[i] = resBytes
	}

	// trigger go routine (report) to start
	if ch, ok := c.ch.rep.Load(vin); ok {
		ch.(repChan) <- res
	}
}

func (c *stubMqttClient) mockStatus(vin int, packet packet) {
	// trigger go routine (status) to start
	if ch, ok := c.ch.sts.Load(vin); ok {
		ch.(stsChan) <- packet
	}
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
