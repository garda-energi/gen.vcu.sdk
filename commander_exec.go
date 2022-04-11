package sdk

import (
	"bytes"
	"reflect"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// exec execute command and return the response.
func (c *commander) exec(invoker string, msg message) (message, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.client.IsConnected() {
		return nil, errClientDisconnected
	}

	cmd, err := getCmdByInvoker(invoker)
	if err != nil {
		return nil, err
	}

	if err := c.sendCommand(cmd, msg); err != nil {
		return nil, err
	}

	return c.waitResponse(cmd)
}

// sendCommand encode and send outgoing command.
func (c *commander) sendCommand(cmd *command, msg message) error {
	if msg.overflow() {
		return errInputOutOfRange("message")
	}

	cp := makeCommandPacket(c.vin, cmd, msg)
	packet, err := encodePacket(cp)
	if err != nil {
		return err
	}

	topic := setTopicVin(TOPIC_COMMAND, c.vin)
	return c.client.pub(topic, 1, true, packet)
}

// waitResponse wait, decode and check of incomming ACK and RESPONSE packet.
func (c *commander) waitResponse(cmd *command) (message, error) {
	defer func() {
		// c.flush()
		c.sleeper.Sleep(3 * time.Second)
	}()

	packet, err := c.waitPacket("ack", DEFAULT_ACK_TIMEOUT)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(packet, strToBytes(PREFIX_ACK)) {
		return nil, errPacketAckCorrupt
	}

	packet, err = c.waitPacket("response", cmd.timeout)
	if err != nil {
		return nil, err
	}

	res, err := decodeResponse(packet)
	if err != nil {
		return nil, err
	}

	if err := res.validateResponse(c.vin, cmd); err != nil {
		return nil, err
	}

	return res.Message, nil
}

// waitPacket wait incomming packet for current VIN.
// It throws error on timeout.
func (c *commander) waitPacket(name string, timeout time.Duration) (packet, error) {
	// flush channel
	for len(c.resChan) > 0 {
		<-c.resChan
	}

	select {
	case data := <-c.resChan:
		return data, nil
	case <-c.sleeper.After(timeout):
		return nil, errPacketTimeout(name)
	}
}

// Destroy unsubscribe from command & response topic for current VIN.
func (c *commander) Destroy() error {
	topics := []string{
		setTopicVin(TOPIC_COMMAND, c.vin),
		setTopicVin(TOPIC_RESPONSE, c.vin),
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
	topic := setTopicVin(TOPIC_COMMAND, c.vin)
	if err := c.client.sub(topic, QOS_SUB_COMMAND, cFunc); err != nil {
		return err
	}

	rFunc := func(client mqtt.Client, msg mqtt.Message) {
		c.logger.Println(CMD, debugPacket(msg))
		c.resChan <- msg.Payload()
	}
	topic = setTopicVin(TOPIC_RESPONSE, c.vin)
	if err := c.client.sub(topic, QOS_SUB_RESPONSE, rFunc); err != nil {
		return err
	}
	return nil
}

// flush clear command & response topic on client
// It indicates that command is done or cancelled.
func (c *commander) flush() {
	for _, t := range []string{TOPIC_COMMAND, TOPIC_RESPONSE} {
		_ = c.client.pub(setTopicVin(t, c.vin), QOS_CMD_FLUSH, true, nil)
	}
}

// invoke call related command using reflection
func (c *commander) invoke(invoker string, arg interface{}) (res, err interface{}) {
	method := reflect.ValueOf(c).MethodByName(invoker)
	ins := []reflect.Value{}
	if arg != nil {
		rv := reflect.ValueOf(arg)
		if (invoker == "McuSpeedMax" || invoker == "McuSetDriveMode") && rv.Kind() == reflect.Slice {
			for i := 0; i < rv.Len(); i++ {
				ins = append(ins, rv.Index(i))
			}
		} else {
			ins = append(ins, rv)
		}
	}
	outs := method.Call(ins)

	err = outs[len(outs)-1].Interface()
	if len(outs) > 1 {
		res = outs[0].Interface()
	}
	return
}
