package sdk

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Destroy unsubscribe from command & response topic for current VIN.
func (c *commander) Destroy() error {
	topics := []string{
		setTopicToVin(TOPIC_COMMAND, c.vin),
		setTopicToVin(TOPIC_RESPONSE, c.vin),
	}
	return c.client.unsub(topics)
}

// listen subscribe to command & response topic for current VIN.
func (c *commander) listen() error {
	if !c.client.IsConnected() {
		return errClientDisconnected
	}

	cFunc := func(client mqtt.Client, msg mqtt.Message) {
		c.logger.Println(CMD, debugPacket(msg))
	}
	topic := setTopicToVin(TOPIC_COMMAND, c.vin)
	if err := c.client.sub(topic, QOS_SUB_COMMAND, cFunc); err != nil {
		return err
	}

	rFunc := func(client mqtt.Client, msg mqtt.Message) {
		c.logger.Println(CMD, debugPacket(msg))
		c.resChan <- msg.Payload()
	}
	topic = setTopicToVin(TOPIC_RESPONSE, c.vin)
	if err := c.client.sub(topic, QOS_SUB_RESPONSE, rFunc); err != nil {
		return err
	}
	return nil
}

// flush clear command & response topic on client
// It indicates that command is done or cancelled.
func (c *commander) flush() {
	for _, t := range []string{TOPIC_COMMAND, TOPIC_RESPONSE} {
		_ = c.client.pub(setTopicToVin(t, c.vin), QOS_CMD_FLUSH, true, nil)
	}
}
