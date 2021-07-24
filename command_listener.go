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
	return c.broker.unsubMulti(topics)
}

// listenResponse subscribe to command & response topic for current VIN.
func (c *commander) listenResponse() error {
	cFunc := func(client mqtt.Client, msg mqtt.Message) {
		if c.logging {
			logPacket(msg)
		}
	}
	if err := c.broker.sub(setTopicToVin(TOPIC_COMMAND, c.vin), QOS_SUB_COMMAND, cFunc); err != nil {
		return err
	}

	rFunc := func(client mqtt.Client, msg mqtt.Message) {
		if c.logging {
			logPacket(msg)
		}
		c.resChan <- msg.Payload()
	}
	if err := c.broker.sub(setTopicToVin(TOPIC_RESPONSE, c.vin), QOS_SUB_RESPONSE, rFunc); err != nil {
		return err
	}
	return nil
}

// flush clear command & response topic on broker.
// It indicates that command is done or cancelled.
func (c *commander) flush() {
	for _, t := range []string{TOPIC_COMMAND, TOPIC_RESPONSE} {
		c.broker.pub(setTopicToVin(t, c.vin), QOS_CMD_FLUSH, true, nil)
	}
}
