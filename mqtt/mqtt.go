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
	Config    ClientConfig
	Listeners Listeners
	client    mqtt.Client
}

func (mq *Mqtt) Connect() error {
	opts := createClientOptions(mq.Config)
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

type Listeners map[string]mqtt.MessageHandler

func (mq *Mqtt) SubscribeAll() error {
	for topic, handler := range mq.Listeners {
		token := mq.client.Subscribe(topic, 1, handler)

		if token.Wait() && token.Error() != nil {
			return token.Error()
		}

		log.Printf("[MQTT] Subscribed to: %s\n", topic)
	}
	return nil
}
