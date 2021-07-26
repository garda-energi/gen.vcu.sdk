package sdk

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// fakeClient implements fake client stub
type fakeClient struct {
	Client
	responses [][]byte
	cmdChan   chan []byte
	resChan   chan struct{}
}

func (b *fakeClient) pub(topic string, qos byte, retained bool, payload []byte) error {
	if flush := payload == nil; !flush {
		b.cmdChan <- payload
	}
	return nil
}

func (b *fakeClient) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	var client mqtt.Client
	msg := &fakeMessage{topic: topic}

	switch toGlobalTopic(topic) {
	case TOPIC_COMMAND:
		go func() {
			select {
			case msg.payload = <-b.cmdChan:
				handler(client, msg)
				b.resChan <- struct{}{}
			case <-time.After(time.Second):
			}
		}()
	case TOPIC_RESPONSE:
		go func() {
			select {
			case <-b.resChan:
				for _, msg.payload = range b.responses {
					time.Sleep(5 * time.Millisecond)
					handler(client, msg)
				}
			case <-time.After(time.Second):
			}
		}()
	}
	return nil
}

func (b *fakeClient) unsub(topics []string) error {
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

func newFakeResponse(vin int, cmdName string) *responsePacket {
	cmd, _ := getCmdByName(cmdName)

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
