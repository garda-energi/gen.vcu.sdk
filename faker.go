package sdk

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// fakeBroker implements fake broker stub
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

func (b *fakeBroker) unsub(topics []string) error {
	return nil
}

// fakeMessage implement fake message stub
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
type fakeSleeper struct{}

func (s *fakeSleeper) Sleep(d time.Duration) {
	d = reduceDuration(d, time.Millisecond)
	time.Sleep(d)
}

func (s *fakeSleeper) After(d time.Duration) <-chan time.Time {
	d = reduceDuration(d, 125*time.Millisecond)
	return time.After(d)
}

// reduceDuration reduce d for faster sleep stub with minimum limit
func reduceDuration(d time.Duration, min time.Duration) time.Duration {
	d /= time.Microsecond
	if d < min {
		d = min
	}
	return d
}
