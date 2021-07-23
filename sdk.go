package gen_vcu_sdk

import (
	cmd "github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
)

type Sdk struct {
	transport *transport.Transport
	logging   bool
}

// New create new instance of Sdk for VCU (Vehicle Control Unit).
func New(host string, port int, user, pass string, logging bool) Sdk {
	tport := transport.New(transport.Config{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	})
	return Sdk{
		transport: tport,
		logging:   logging,
	}
}

// Connect open connection to mqtt broker.
func (s *Sdk) Connect() error {
	return s.transport.Connect()
}

// Disconnect close connection to mqtt broker.
func (s *Sdk) Disconnect() {
	s.transport.Disconnect()
}

// NewCommand create new instance of Command for specific VIN.cmdVins
func (s *Sdk) NewCommand(vin int) (*cmd.Command, error) {
	return cmd.New(vin, s.transport)
}

// AddListener subscribe to Status & Report topic (if callback is specified) for spesific vin in range.
//
// Examples :
//
// listen by list :
// s.AddListener([]int{1, 2 ,3}, *listerner)
//
// listen one spesific vin 2341 :
// s.AddListener([]int{2341}, *listerner)
//
// listen by range :
// s.AddListener(sdk.VinRange(min, max), *listerner)
func (s *Sdk) AddListener(vins []int, l *Listener) error {
	if l.StatusFunc != nil {
		topic := setTopicToVins(shared.TOPIC_STATUS, vins)
		if err := s.transport.SubMulti(topic, 1, StatusListener(l.StatusFunc, s.logging)); err != nil {
			return err
		}
	}

	if l.ReportFunc != nil {
		topic := setTopicToVins(shared.TOPIC_REPORT, vins)
		if err := s.transport.SubMulti(topic, 1, ReportListener(l.ReportFunc, s.logging)); err != nil {
			return err
		}
	}
	return nil
}

// RemoveListener unsubscribe status topic and report for spesific vin in range.
func (s *Sdk) RemoveListener(vins []int) error {
	// topics is status topic + report topic
	topics := append(setTopicToVins(shared.TOPIC_STATUS, vins), setTopicToVins(shared.TOPIC_REPORT, vins)...)
	return s.transport.UnsubMulti(topics)
}

// VinRange will generate array of integer from min to max.
//
// If min greater than max, it will be switched
func VinRange(min int, max int) []int {
	// switch them if min greater than max
	if max < min {
		tmpMin := min
		min = max
		max = tmpMin
	}
	// generate sequence number
	len := max - min + 1
	result := make([]int, len)
	for i := range result {
		result[i] = min + i
	}
	return result
}

// setTopicToVins create multiple topic for list of vin
func setTopicToVins(topic string, vins []int) []string {
	topics := make([]string, len(vins))
	for i, v := range vins {
		topics[i] = shared.SetTopicToVin(topic, v)
	}
	return topics
}
