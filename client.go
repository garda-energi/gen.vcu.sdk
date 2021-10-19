package sdk

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type subscriber struct {
	qos     byte
	handler mqtt.MessageHandler
}

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
	logger      *log.Logger
	subscribers *sync.Map
}

func init() {
	// mqtt.DEBUG = log.New(os.Stderr, "DEBUG - ", log.LstdFlags)
	// mqtt.CRITICAL = log.New(os.Stderr, "CRITICAL - ", log.LstdFlags)
	// mqtt.WARN = log.New(os.Stderr, "WARN - ", log.LstdFlags)
	mqtt.ERROR = log.New(os.Stderr, "ERROR - ", log.LstdFlags)
}

// newClient create instance of mqtt client
func newClient(config *ClientConfig, logger *log.Logger) *client {
	client := client{
		logger:      logger,
		subscribers: &sync.Map{},
	}
	client.Client = mqtt.NewClient(client.newClientOptions(config))
	return &client
}

func newFakeClient(fakeClient mqtt.Client, logger *log.Logger) *client {
	client := client{
		logger:      logger,
		subscribers: &sync.Map{},
		Client:      fakeClient,
	}
	_ = client.newClientOptions(&ClientConfig{})
	return &client
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

	c.subscribers.Store(topic, subscriber{qos: qos, handler: handler})
	c.logger.Println(CLI, "Subscribed to:", topic)
	return nil
}

func (c *client) subMulti(topics []string, qos byte, handler mqtt.MessageHandler) error {
	topicFilters := make(map[string]byte, len(topics))
	for _, topic := range topics {
		topicFilters[topic] = qos
	}

	token := c.SubscribeMultiple(topicFilters, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, topic := range topics {
		c.subscribers.Store(topic, subscriber{qos: qos, handler: handler})
		c.logger.Println(CLI, "Subscribed to:", topic)
	}
	return nil
}

func (c *client) unsub(topics []string) error {
	token := c.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, topic := range topics {
		c.subscribers.Delete(topic)
		c.logger.Println(CLI, "Un-subscribed from:", topic)
	}
	return nil
}

func (c *client) newClientOptions(cfg *ClientConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port))
	opts.SetClientID(fmt.Sprintf("go_mqtt_client_%d", time.Now().Unix()))
	opts.SetUsername(cfg.User)
	opts.SetPassword(cfg.Pass)
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		c.logger.Println(CLI, debugPacket(msg))
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		c.logger.Println(CLI, "Connected")

		// resubscribe all topics or reconnect
		c.subscribers.Range(func(key, value interface{}) bool {
			topic := key.(string)
			subs := value.(subscriber)

			err := c.sub(topic, subs.qos, subs.handler)
			return err == nil
		})
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		c.logger.Println(CLI, "Disconnected", err)
	})

	return opts
}
