package sdk

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ClientConfig store connection string for mqtt client
type ClientConfig struct {
	Host string
	Port int
	User string
	Pass string
}

// client implements mqtt client
type client struct {
	mqtt.Client
	logger *log.Logger
}

// newClient create instance of mqtt client
func newClient(config *ClientConfig, logger *log.Logger) *client {
	opts := newClientOptions(config, logger)
	return &client{
		Client: mqtt.NewClient(opts),
		logger: logger,
	}
}

func (c *client) pub(topic string, qos byte, retained bool, packet packet) error {
	token := c.Publish(topic, qos, retained, []byte(packet))
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	c.logger.Println(CLI, "Published to:", topic)
	return nil
}

func (c *client) sub(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := c.Subscribe(topic, qos, handler)

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	c.logger.Println(CLI, "Subscribed to:", topic)
	return nil
}

func (c *client) subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	topicFilters := make(map[string]byte, len(topics))
	for _, v := range topics {
		topicFilters[v] = qos
	}

	token := c.SubscribeMultiple(topicFilters, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, v := range topics {
		c.logger.Println(CLI, "Subscribed to:", v)
	}
	return nil
}

func (c *client) unsub(topics []string) error {
	token := c.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, v := range topics {
		c.logger.Println(CLI, "Un-subscribed from:", v)
	}
	return nil
}

// newClientOptions make client options for mqtt.
func newClientOptions(c *ClientConfig, logger *log.Logger) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", c.Host, c.Port))
	opts.SetClientID(fmt.Sprintf("go_mqtt_client_%d", time.Now().Unix()))
	opts.SetUsername(c.User)
	opts.SetPassword(c.Pass)
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		logger.Println(CLI, debugPacket(msg))
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logger.Println(CLI, "Connected")
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logger.Println(CLI, "Disconnected", err)
	})

	return opts
}
