package api

import (
	"encoding/hex"
	"log"
	"strings"

	// ttt "github.com/eclipse/paho.transport.golang"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

const (
	TOPIC_REPORT = "VCU/+/RPT"
)

type CallbackFunc func(interface{}, []byte)

type Api struct {
	config         transport.ClientConfig
	reportCallback CallbackFunc
}

func New(host string, port int, user, pass string) Api {
	return Api{
		config: transport.ClientConfig{
			Host:     host,
			Port:     port,
			Username: user,
			Password: pass,
			ClientId: "go_mqtt_client",
		},
	}
}

func (a *Api) AddReportListener(cb CallbackFunc) {
	a.reportCallback = cb
}

func (a *Api) ConnectAndListen() {
	t := transport.New(a.config)

	if err := t.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %s\n", err.Error())
	}

	if err := t.Subscribe(TOPIC_REPORT, a.reportHandler); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %s\n", err.Error())
	}

	util.WaitForCtrlC()
	t.Disconnect()
}

func (a *Api) reportHandler(client mqtt.Client, msg mqtt.Message) {
	hexString := strings.ToUpper(hex.EncodeToString(msg.Payload()))
	log.Printf("[REPORT] %s\n", hexString)

	packet := &Report{Bytes: msg.Payload()}
	report, err := packet.decodeReport()
	if err != nil {
		log.Fatal(err)
	}
	// util.Debug(report)
	a.reportCallback(report, msg.Payload())
}
