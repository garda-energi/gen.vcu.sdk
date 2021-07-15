package transport

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type clientConfig struct {
	Host     string
	Port     int
	ClientId string
	Username string
	Password string
}

type Transport struct {
	config clientConfig
	client mqtt.Client
}

func New(host string, port int, user, pass string) Transport {
	return Transport{config: clientConfig{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		// ClientId: "go_mqtt_client",
	}}
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
