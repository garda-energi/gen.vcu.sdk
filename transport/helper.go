package transport

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func createClientOptions(config TransportConfig) *mqtt.ClientOptions {
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

func defaultPublishHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] Topic: %s => %s\n", msg.Topic(), msg.Payload())
}

func connectHandler(client mqtt.Client) {
	log.Printf("[MQTT] Connected\n")
}

func disconnectHandler(client mqtt.Client, err error) {
	log.Printf("[MQTT] Disconnected, %v\n", err)
}
