package broker

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Host string
	Port int
	User string
	Pass string
}

type Broker struct {
	config Config
	client mqtt.Client
}

// New create instance of Broker.
func New(config Config) *Broker {
	return &Broker{config: config}
}

// Connect open connection to mqtt broker.
func (b *Broker) Connect() error {
	opts := createClientOptions(b.config)
	b.client = mqtt.NewClient(opts)

	token := b.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Disconnect close connection to mqtt broker.
func (b *Broker) Disconnect() {
	b.client.Disconnect(100)
}

// Sub subscribe to mqtt topic.
func (b *Broker) Sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := b.client.Subscribe(topic, qos, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] Subscribed to: %s\n", topic)
	return nil
}

// SubMulti subscribe to muliple mqtt topics.
func (b *Broker) SubMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	topicFilters := map[string]byte{}
	for _, v := range topics {
		topicFilters[v] = qos
	}

	token := b.client.SubscribeMultiple(topicFilters, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, v := range topics {
		log.Printf("[MQTT] Subscribed to: %s\n", v)
	}
	return nil
}

// UnsubMulti unsubscribe mqtt muliple topic.
func (b *Broker) UnsubMulti(topics []string) error {
	token := b.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Pub publish to mqtt topic.
func (b *Broker) Pub(topic string, qos byte, retained bool, payload []byte) {
	token := b.client.Publish(topic, qos, retained, payload)
	token.Wait()

	log.Printf("[MQTT] Published to: %s\n", topic)
}
