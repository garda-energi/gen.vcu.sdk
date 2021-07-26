package sdk

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client is building block for client client (with extra things).
type Client interface {
	mqtt.Client
	// pub publish to mqtt topic.
	pub(topic string, qos byte, retained bool, payload []byte) error
	// sub subscribe to mqtt topic.
	sub(topic string, qos byte, handler mqtt.MessageHandler) error
	// subMulti subscribe to mqtt topics.
	subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error
	// unsub unsubscribe from mqtt topics.
	unsub(topics []string) error
}

// ClientConfig store connection string for client client
type ClientConfig struct {
	Host string
	Port int
	User string
	Pass string
}
type client struct {
	mqtt.Client
	logger *log.Logger
}

// newClient create instance of Client client.
func newClient(config *ClientConfig, logging bool) Client {
	logger := newLogger(logging, "CLIENT")
	return &client{
		Client: mqtt.NewClient(newClientOptions(config, logger)),
		logger: logger,
	}
}

func (b *client) pub(topic string, qos byte, retained bool, payload []byte) error {
	token := b.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	b.logger.Printf("Published to: %s\n", topic)
	return nil
}

func (b *client) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := b.Subscribe(topic, qos, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	b.logger.Printf("Subscribed to: %s\n", topic)
	return nil
}

func (b *client) subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	topicFilters := map[string]byte{}
	for _, v := range topics {
		topicFilters[v] = qos
	}

	token := b.SubscribeMultiple(topicFilters, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, v := range topics {
		b.logger.Printf("Subscribed to: %s\n", v)
	}
	return nil
}

func (b *client) unsub(topics []string) error {
	token := b.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, v := range topics {
		b.logger.Printf("Un-subscribed from: %s\n", v)
	}
	return nil
}

// newClientOptions make client options for mqtt.
func newClientOptions(c *ClientConfig, logger *log.Logger) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.Host, c.Port))
	opts.SetUsername(c.User)
	opts.SetPassword(c.Pass)
	opts.SetClientID("go_mqtt_client")

	opts.DefaultPublishHandler = func(client mqtt.Client, msg mqtt.Message) {
		logger.Println(debugPacket(msg))
	}
	opts.OnConnect = func(client mqtt.Client) {
		logger.Println("Connected")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logger.Printf("Disconnected, %v\n", err)
	}
	return opts
}
