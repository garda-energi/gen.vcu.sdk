package command

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

// Listen subscribe to Command & Response topic for multiple vins.
func (c *Command) listen() error {
	cFunc := func(client mqtt.Client, msg mqtt.Message) {
		shared.LogMessage(msg)
		// vin := util.TopicVin(msg.Topic())
	}
	topic := shared.SetTopicToVin(shared.TOPIC_COMMAND, c.vin)
	if err := c.transport.Sub(topic, 1, cFunc); err != nil {
		return err
	}

	rFunc := func(client mqtt.Client, msg mqtt.Message) {
		shared.LogMessage(msg)
		// vin := shared.GetTopicVin(msg.Topic())
		// responses.set(vin, msg.Payload())
		c.resChan <- msg.Payload()
	}
	topic = shared.SetTopicToVin(shared.TOPIC_RESPONSE, c.vin)
	if err := c.transport.Sub(topic, 1, rFunc); err != nil {
		return err
	}
	return nil
}

// Destroy unsubscribe status topic and report for spesific vin in range.
func (c *Command) Destroy() error {
	topics := []string{
		shared.SetTopicToVin(shared.TOPIC_COMMAND, c.vin),
		shared.SetTopicToVin(shared.TOPIC_RESPONSE, c.vin),
	}
	return c.transport.UnsubMulti(topics)
}
