package mqtt

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

type Mqtt struct {
	client mqtt.Client
}

func (mq *Mqtt) Connect() error {
	opts := createClientOptions(ClientConfig{
		Host:     "test.mosquitto.org",
		Port:     1883,
		ClientId: "go_mqtt_client",
		Username: "",
		Password: "",
	})
	mq.client = mqtt.NewClient(opts)

	token := mq.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mq *Mqtt) Disconnect() {
	mq.client.Disconnect(100)
}

type Subscribers map[string]mqtt.MessageHandler

func (mq *Mqtt) Subscribe(subscribers Subscribers) error {
	for topic, handler := range subscribers {
		token := mq.client.Subscribe(topic, 1, handler)

		if token.Wait() && token.Error() != nil {
			return token.Error()
		}

		log.Printf("[MQTT] Subscribed to: %s\n", topic)
	}
	return nil
}
