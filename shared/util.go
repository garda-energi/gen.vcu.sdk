package shared

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

func Byte2Hex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

func Hex2Byte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func LogMessage(msg mqtt.Message) {
	log.Printf("[%s] %s\n", msg.Topic(), Byte2Hex(msg.Payload()))
}

func Reverse(b []byte) []byte {
	nb := make([]byte, len(b))
	for i := range nb {
		nb[i] = b[len(b)-1-i]
	}
	return nb
}

func GetTopicVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])
	return vin
}

// SetTopicToVin change willcard topic to spesific topic
func SetTopicToVin(topic string, vin int) string {
	return strings.Replace(topic, "+", strconv.Itoa(vin), 1)
}
