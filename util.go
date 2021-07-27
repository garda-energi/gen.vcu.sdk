package sdk

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
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
	return log.New(out, fmt.Sprint(prefix, " "), log.Ldate|log.Ltime)
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

func getPacketSize(v interface{}) int {
	size := 0
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}

	switch rk := rv.Kind(); rk {

	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			rvField := rv.Field(i)
			rtField := rv.Type().Field(i)

			tagField := deTag(rtField.Tag, rvField.Kind())
			if rvField.Type() == typeOfTime {
				size += tagField.Len
			} else if rk := rvField.Kind(); rk == reflect.Struct || rk == reflect.Array || rk == reflect.Ptr || rk == reflect.Slice {
				size += getPacketSize(rvField.Addr().Interface())
			} else {
				size += tagField.Len
			}
		}

	case reflect.Array, reflect.Slice:
		if rv.Type() == typeOfMessage {
			size += rv.Len()
		} else {
			for i := 0; i < rv.Len(); i++ {
				size += getPacketSize(rv.Index(i).Addr().Interface())
			}
		}

	default:
	}

	return size
}
