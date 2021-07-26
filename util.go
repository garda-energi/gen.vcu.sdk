package sdk

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SetupGracefulShutdown wait until ctrl+c is pressed
func SetupGracefulShutdown() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		fmt.Println("Gracefully exit application")
		os.Exit(1)
	}()

	return stopChan
}

func newLogger(logging bool, prefix string) *log.Logger {
	out := ioutil.Discard
	if logging {
		out = os.Stderr
	}
	return log.New(out, fmt.Sprintf("[%s] ", prefix), log.Ldate|log.Ltime)
}

// // randomSleep will sleep betwen random time specified
// func randomSleep(min, max int, unit time.Duration) {
// 	rng := rand.Intn(max-min) + min
// 	time.Sleep(time.Duration(rng) * unit)
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

// debugPacket format received mqtt message
func debugPacket(msg mqtt.Message) string {
	return fmt.Sprintf("%s => %s\n", msg.Topic(), byteToHex(msg.Payload()))
}

// reverseBytes swap bytes position
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

// toGlobalTopic replace VIN from topic with '+' character
func toGlobalTopic(topic string) string {
	s := strings.Split(topic, "/")
	s[1] = "+"
	return strings.Join(s, "/")
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
