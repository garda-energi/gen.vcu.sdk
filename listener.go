package gen_vcu_sdk

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type StatusListenerFunc func(vin int, online bool) error
type ReportListenerFunc func(vin int, report *report.ReportPacket) error

type Listener struct {
	StatusFunc StatusListenerFunc
	ReportFunc ReportListenerFunc
}

func listenerHook(msg mqtt.Message, logging bool) int {
        if logging {
                util.LogMessage(msg)                                                                         }
        vin := util.TopicVin(msg.Topic())
        return vin
}

func parseOnline(b []byte) bool {
	return b[0] == '1'
}
