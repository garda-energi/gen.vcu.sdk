package handler

import (
	"bytes"
	"encoding/hex"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen-go-packet/decoder"
)

func Report(client mqtt.Client, msg mqtt.Message) {
	hexString := strings.ToUpper(hex.EncodeToString(msg.Payload()))
	log.Printf("[REPORT] %s\n", hexString)

	reader := bytes.NewReader(msg.Payload())
	report := &decoder.Report{Reader: reader}

	_, err := report.Decode()
	if err != nil {
		log.Fatal(err)
	}

	// util.Debug(reportDecoded)
}
