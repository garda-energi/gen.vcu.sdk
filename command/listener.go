package command

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

// ResponseListener executed when got new packet on response topic
func ResponseListener(client mqtt.Client, msg mqtt.Message) {
	vin := util.TopicVin(msg.Topic())

	responses.set(vin, msg.Payload())

	util.LogMessage(msg)
}

// CommandListener executed when got new packet on command topic
func CommandListener(client mqtt.Client, msg mqtt.Message) {
	// vin := util.TopicVin(msg.Topic())

	util.LogMessage(msg)
}
