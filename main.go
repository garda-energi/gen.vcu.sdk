package main

import (
	"log"

	"github.com/pudjamansyurin/gen-go-packet/handler"
	"github.com/pudjamansyurin/gen-go-packet/mqtt"
	"github.com/pudjamansyurin/gen-go-packet/util"
)

func main() {
	mq := &mqtt.Mqtt{}
	if err := mq.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %s\n", err.Error())
	}

	subscribers := mqtt.Subscribers{
		"VCU/+/RPT": handler.Report,
	}
	if err := mq.Subscribe(subscribers); err != nil {
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
