package transport

import (
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type Config struct {
	Host string
	Port int
	User string
	Pass string
}

type Transport struct {
	config Config
	client mqtt.Client
}

// New create instance of Transport.
func New(config Config) *Transport {
	return &Transport{config: config}
}

// Connect open connection to mqtt broker.
func (t *Transport) Connect() error {
	opts := createClientOptions(t.config)
	t.client = mqtt.NewClient(opts)

	token := t.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Disconnect close connection to mqtt broker.
func (t *Transport) Disconnect() {
	t.client.Disconnect(100)
}

// Sub subscribe to mqtt topic.
func (t *Transport) Sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := t.client.Subscribe(topic, qos, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] Subscribed to: %s\n", topic)
	return nil
}

// SubMulti subscribe to muliple mqtt topics.
func (t *Transport) SubMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	topicFilters := map[string]byte{}
	vins := make([]string, len(topics))
	for i, v := range topics {
		topicFilters[v] = qos
		vins[i] = strconv.Itoa(shared.GetTopicVin(v))
	}

	token := t.client.SubscribeMultiple(topicFilters, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] Subscribed to: %v\n", strings.Join(vins, ", "))
	return nil
}

// UnsubMulti unsubscribe mqtt muliple topic.
func (t *Transport) UnsubMulti(topics []string) error {
	token := t.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Pub publish to mqtt topic.
func (t *Transport) Pub(topic string, qos byte, retained bool, payload []byte) {
	token := t.client.Publish(topic, qos, retained, payload)
	token.Wait()

	log.Printf("[MQTT] Published to: %s\n", topic)
}
