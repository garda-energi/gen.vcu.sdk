package transport

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TransportConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type Transport struct {
	config TransportConfig
	Client mqtt.Client
}

func New(config TransportConfig) Transport {
	return Transport{config: config}
}

func (t *Transport) Connect() error {
	opts := createClientOptions(t.config)
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	t.Client = client
	return nil
}

func (t *Transport) Disconnect() {
	t.Client.Disconnect(100)
}

func (t *Transport) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := t.Client.Subscribe(topic, 1, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] Subscribed to: %s\n", topic)
	return nil
}
