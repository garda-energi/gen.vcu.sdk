package sdk

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// listen subscribe to command & response topic for current VIN.
func (c *commander) listen() error {
	cFunc := func(client mqtt.Client, msg mqtt.Message) {
		logPacket(msg)
	}
	if err := c.broker.Sub(setTopicToVin(TOPIC_COMMAND, c.vin), 1, cFunc); err != nil {
		return err
	}

	rFunc := func(client mqtt.Client, msg mqtt.Message) {
		logPacket(msg)
		c.resChan <- msg.Payload()
	}
	if err := c.broker.Sub(setTopicToVin(TOPIC_RESPONSE, c.vin), 1, rFunc); err != nil {
		return err
	}
	return nil
}

// Destroy unsubscribe from command & response topic for current VIN.
func (c *commander) Destroy() error {
	topics := []string{
		setTopicToVin(TOPIC_COMMAND, c.vin),
		setTopicToVin(TOPIC_RESPONSE, c.vin),
	}
	return c.broker.UnsubMulti(topics)
}
