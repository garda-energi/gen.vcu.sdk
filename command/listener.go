package command

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func ResponseListener(client mqtt.Client, msg mqtt.Message) {
	vin := util.TopicVin(msg.Topic())

	responses.set(vin, msg.Payload())

	util.LogMessage(msg)
}

func CommandListener(client mqtt.Client, msg mqtt.Message) {
	// vin := util.TopicVin(msg.Topic())

	util.LogMessage(msg)
}
