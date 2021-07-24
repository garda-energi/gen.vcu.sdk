package sdk

import (
	"encoding/hex"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// // waitForCtrlC wait until ctrl+c is pressed
// func waitForCtrlC() {
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt)
// 	<-stop
// }

// // dd print detailed variable information
// func dd(data interface{}) {
// 	fmt.Printf("%+v\n", data)
// }

// byteToHex convert bytes to hex string
func byteToHex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

// hexToByte convert hex string to bytes
func hexToByte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

// logPacket print received mqtt message
func logPacket(msg mqtt.Message) {
	log.Printf("[%s] %s\n", msg.Topic(), byteToHex(msg.Payload()))
}

// Reverse swap bytes position
func reverseBytes(b []byte) []byte {
	nb := make([]byte, len(b))
	for i := range nb {
		nb[i] = b[len(b)-1-i]
	}
	return nb
}

// getTopicVin extract VIN information from mqtt topic
func getTopicVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])
	return vin
}

// setTopicToVin insert VIN into topic pattern
func setTopicToVin(topic string, vin int) string {
	return strings.Replace(topic, "+", strconv.Itoa(vin), 1)
}

// setTopicToVins create multiple topic for list of vin
func setTopicToVins(topic string, vins []int) []string {
	topics := make([]string, len(vins))
	for i, v := range vins {
		topics[i] = setTopicToVin(topic, v)
	}
	return topics
}
