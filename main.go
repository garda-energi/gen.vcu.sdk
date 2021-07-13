package main

import (
	"log"

	"github.com/pudjamansyurin/gen-go-packet/handler"
	"github.com/pudjamansyurin/gen-go-packet/mqtt"
	"github.com/pudjamansyurin/gen-go-packet/util"
)

func main() {
	mq := &mqtt.Mqtt{
		Config: mqtt.ClientConfig{
			Host:     "test.mosquitto.org",
			Port:     1883,
			ClientId: "go_mqtt_client",
		},
		Listeners: mqtt.Listeners{
			"VCU/+/RPT": handler.Report,
		},
	}

	if err := mq.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %s\n", err.Error())
	}

	if err := mq.SubscribeAll(); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %s\n", err.Error())
	}

	// num := 10
	// for i := 0; i < num; i++ {
	// 	text := fmt.Sprintf("%d", i)
	// 	token = client.Publish(topic, 0, false, text)
	// 	token.Wait()
	// 	time.Sleep(time.Second)
	// }

	// gracefully quit
	util.WaitForCtrlC()
	mq.Disconnect()
}
