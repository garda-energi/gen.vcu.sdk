package handler

import (
	"bytes"
	"encoding/hex"
	"log"
	"reflect"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen-go-packet/decoder"
	"github.com/pudjamansyurin/gen-go-packet/model"
	"github.com/pudjamansyurin/gen-go-packet/util"
)

func ReportHandler(client mqtt.Client, msg mqtt.Message) {
	hexString := strings.ToUpper(hex.EncodeToString(msg.Payload()))
	log.Printf("[REPORT] %s\n", hexString)

	packet := &Report{Bytes: msg.Payload()}
	report, err := packet.decodeReport()
	if err != nil {
		log.Fatal(err)
	}

	util.Debug(report)
}

type Report struct {
	Bytes []byte
}

func (r *Report) decodeReport() (interface{}, error) {
	header := model.HeaderReport{}
	if err := r.decode(&header); err != nil {
		return nil, err
	}

	var report interface{}
	if r.simpleFrame(header) {
		simple := model.ReportSimple{}
		if err := r.decode(&simple); err != nil {
			return nil, err
		}
		report = simple
	} else {
		full := model.ReportFull{}
		if err := r.decode(&full); err != nil {
			return nil, err
		}
		report = full
	}

	return report, nil
}

func (r *Report) decode(dst interface{}) error {
	return decoder.TagWalk(bytes.NewReader(r.Bytes), reflect.ValueOf(dst), "")
}

func (r *Report) simpleFrame(header model.HeaderReport) bool {
	return header.FrameID == model.FRAME_ID_SIMPLE
}
