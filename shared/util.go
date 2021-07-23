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

// WaitForCtrlC wait until ctrl+c is pressed
func WaitForCtrlC() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

// Debug print detailed variable information
func Debug(data interface{}) {
	fmt.Printf("%+v\n", data)
}

// Byte2Hex convert bytes to hex string
func Byte2Hex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

// Hex2Byte convert hex string to bytes
func Hex2Byte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

// LogMessage print received mqtt message
func LogMessage(msg mqtt.Message) {
	log.Printf("[%s] %s\n", msg.Topic(), Byte2Hex(msg.Payload()))
}

// Reverse swap bytes position
func Reverse(b []byte) []byte {
	nb := make([]byte, len(b))
	for i := range nb {
		nb[i] = b[len(b)-1-i]
	}
	return nb
}

// GetTopicVin extract VIN information from mqtt topic
func GetTopicVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])
	return vin
}

// SetTopicToVin insert VIN into topic pattern
func SetTopicToVin(topic string, vin int) string {
	return strings.Replace(topic, "+", strconv.Itoa(vin), 1)
}
