package util

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func WaitForCtrlC() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func Debug(data interface{}) {
	fmt.Printf("%+v\n", data)
	// spew.Dump(data)
}

func HexString(payload []byte) string {
	return strings.ToUpper(hex.EncodeToString(payload))
}

func LogMessage(msg mqtt.Message) {
	log.Printf("[%s] %s\n", msg.Topic(), HexString(msg.Payload()))
}

func Reverse(b []byte) []byte {
	nb := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		nb[i] = b[len(b)-1-i]
	}
	return nb
}

func TopicVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])

	return vin
}
