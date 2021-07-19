package command

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func (c *Command) ResponseListener(client mqtt.Client, msg mqtt.Message) {
	vin := util.TopicVin(msg.Topic())

	RX.Set(vin, msg.Payload())

	util.LogMessage(msg)
}

func (c *Command) CommandListener(client mqtt.Client, msg mqtt.Message) {
	// vin := util.TopicVin(msg.Topic())

	util.LogMessage(msg)
}
