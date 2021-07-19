package transport

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

type Transport struct {
	config Config
	Client mqtt.Client
}

func New(config Config) *Transport {
	return &Transport{config: config}
}

func (t *Transport) Connect() error {
	opts := createClientOptions(t.config)
	t.Client = mqtt.NewClient(opts)

	token := t.Client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

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
