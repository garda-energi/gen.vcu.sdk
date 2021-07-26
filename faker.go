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
type fakeSleeper struct{}

func (s *fakeSleeper) Sleep(d time.Duration) {
	d = faster(d, time.Millisecond)
	time.Sleep(d)
}

func (s *fakeSleeper) After(d time.Duration) <-chan time.Time {
	d = faster(d, 125*time.Millisecond)
	return time.After(d)
}

// faster reduce d for faster sleep stub with minimum duration
func faster(d time.Duration, min time.Duration) time.Duration {
	d /= time.Microsecond
	if d < min {
		d = min
	}
	return d
}

func newFakeResponse(vin int, cmdName string) *responsePacket {
	cmd, _ := getCommand(cmdName)

	return &responsePacket{
		Header: &headerResponse{
			HeaderCommand: HeaderCommand{
				Header: Header{
					Prefix:       PREFIX_RESPONSE,
					Size:         0,
					Vin:          uint32(vin),
					SendDatetime: time.Now(),
				},
				Code:    cmd.code,
				SubCode: cmd.subCode,
			},
			ResCode: resCodeOk,
		},
		Message: nil,
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
