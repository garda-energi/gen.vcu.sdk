package sdk

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type brokerConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type Broker interface {
	// connect open connection to mqtt broker.
	connect() error
	// disconnect close connection to mqtt broker.
	disconnect()
	// pub publish to mqtt topic.
	pub(topic string, qos byte, retained bool, payload []byte) error
	// sub subscribe to mqtt topic.
	sub(topic string, qos byte, handler mqtt.MessageHandler) error
	// subMulti subscribe to muliple mqtt topics.
	subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error
	// unsubMulti unsubscribe mqtt muliple topic.
	unsubMulti(topics []string) error
}

type broker struct {
	config brokerConfig
	client mqtt.Client
}

// newBroker create instance of Broker.
func newBroker(config brokerConfig) Broker {
	return &broker{config: config}
}

// connect open connection to mqtt broker.
func (b *broker) connect() error {
	opts := newClientOptions(b.config)
	b.client = mqtt.NewClient(opts)

	token := b.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// disconnect close connection to mqtt broker.
func (b *broker) disconnect() {
	b.client.Disconnect(100)
}

// pub publish to mqtt topic.
func (b *broker) pub(topic string, qos byte, retained bool, payload []byte) error {
	token := b.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("[MQTT] Published to: %s\n", topic)
	return nil
}

// sub subscribe to mqtt topic.
func (b *broker) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := b.client.Subscribe(topic, qos, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] Subscribed to: %s\n", topic)
	return nil
}

// subMulti subscribe to muliple mqtt topics.
func (b *broker) subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
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

// unsubMulti unsubscribe mqtt muliple topic.
func (b *broker) unsubMulti(topics []string) error {
	token := b.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, v := range topics {
		log.Printf("[MQTT] Un-subscribed from: %s\n", v)
	}
	return nil
}

// newClientOptions make client options for mqtt.
func newClientOptions(config brokerConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Host, config.Port))
	opts.SetUsername(config.User)
	opts.SetPassword(config.Pass)
	opts.SetClientID("go_mqtt_client")

	opts.SetDefaultPublishHandler(defaultPublishHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = disconnectHandler

	return opts
}

// defaultPublishHandler executed when no publish handler specified.
func defaultPublishHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] Topic: %s => %s\n", msg.Topic(), msg.Payload())
}

// connectHandler executed when mqtt connection is ready.
func connectHandler(client mqtt.Client) {
	log.Printf("[MQTT] Connected\n")
}

// disconnectHandler executed when mqtt is disconnected.
func disconnectHandler(client mqtt.Client, err error) {
	log.Printf("[MQTT] Disconnected, %v\n", err)
}