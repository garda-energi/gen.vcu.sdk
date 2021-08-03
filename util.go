package sdk

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
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

// newLogger create new logger
func newLogger(logging bool, prefix string) *log.Logger {
	out := ioutil.Discard
	if logging {
		out = os.Stderr
	}
	return log.New(out, fmt.Sprint(prefix, " "), log.Ltime)
}

// byteToHex convert bytes to hex string
func byteToHex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

// hexToByte convert hex string to bytes
func hexToByte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

// bitSet check is bit index is set on word
func bitSet(word uint32, bit uint8) bool {
	return word&(1<<bit) > 0
}

// sliceToStr convert sclice to string
func sliceToStr(s interface{}, prefix string) string {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Slice {
		return ""
	}

	buf := make([]string, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		buf[i] = fmt.Sprintf("%s", rv.Index(i))
	}
	return prefix + "[" + strings.Join(buf, ", ") + "]"
}

// debugPacket format received mqtt message
func debugPacket(msg mqtt.Message) string {
	return fmt.Sprintln(msg.Topic(), "=>", byteToHex(msg.Payload()))
}

// reverseBytes swap bytes position
func reverseBytes(b []byte) []byte {
	nb := make([]byte, len(b))
	for i := range nb {
		nb[i] = b[len(b)-1-i]
	}
	return nb
}

// toGlobalTopic replace VIN from topic with '+' character
func toGlobalTopic(topic string) string {
	s := strings.Split(topic, "/")
	s[1] = "+"
	return strings.Join(s, "/")
}

// getTopicVin extract VIN information from mqtt topic
func getTopicVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])
	return vin
}

// setTopicVin insert VIN into topic pattern
func setTopicVin(topic string, vin int) string {
	return strings.Replace(topic, "+", strconv.Itoa(vin), 1)
}

// setTopicVins create multiple topic for list of vin
func setTopicVins(topic string, vins []int) []string {
	topics := make([]string, len(vins))
	for i, v := range vins {
		topics[i] = setTopicVin(topic, v)
	}
	return topics
}

// randBool generate random boolean
func randBool() bool {
	return rand.Intn(1) == 1
}

// randInt generate random integer
func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// randFloat generate random float
func randFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
