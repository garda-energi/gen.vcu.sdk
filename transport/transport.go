package transport

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ClientConfig struct {
	Host     string
	Port     int
	ClientId string
	Username string
	Password string
}

type Transport struct {
	config ClientConfig
	client mqtt.Client
}

func New(config ClientConfig) Transport {
	return Transport{config: config}
}

func (t *Transport) Connect() error {
	opts := createClientOptions(t.config)
	t.client = mqtt.NewClient(opts)

	token := t.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Transport) Disconnect() {
	t.client.Disconnect(100)
}

func (t *Transport) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := t.client.Subscribe(topic, 1, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] Subscribed to: %s\n", topic)
	return nil
}
