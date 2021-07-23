package gen_vcu_sdk

import (
	"github.com/pudjamansyurin/gen_vcu_sdk/broker"
	cmd "github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type Sdk struct {
	broker  *broker.Broker
	logging bool
}

// New create new instance of Sdk for VCU (Vehicle Control Unit).
func New(host string, port int, user, pass string, logging bool) Sdk {
	broker := broker.New(broker.Config{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	})
	return Sdk{
		broker:  broker,
		logging: logging,
	}
}

// Connect open connection to mqtt broker.
func (s *Sdk) Connect() error {
	return s.broker.Connect()
}

// Disconnect close connection to mqtt broker.
func (s *Sdk) Disconnect() {
	s.broker.Disconnect()
}

// NewCommander create new instance of Commander for specific VIN.
func (s *Sdk) NewCommander(vin int) (*cmd.Commander, error) {
	return cmd.New(vin, s.broker)
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
		if err := s.broker.SubMulti(topic, 1, StatusListener(l.StatusFunc, s.logging)); err != nil {
			return err
		}
	}

	if l.ReportFunc != nil {
		topic := setTopicToVins(shared.TOPIC_REPORT, vins)
		if err := s.broker.SubMulti(topic, 1, ReportListener(l.ReportFunc, s.logging)); err != nil {
			return err
		}
	}
	return nil
}

// RemoveListener unsubscribe status and report topic for spesific vin in range.
func (s *Sdk) RemoveListener(vins []int) error {
	topics := append(
		setTopicToVins(shared.TOPIC_STATUS, vins),
		setTopicToVins(shared.TOPIC_REPORT, vins)...,
	)
	return s.broker.UnsubMulti(topics)
}

// VinRange generate array of integer from min to max.
//
// If min greater than max, it will be swapped.
func VinRange(min int, max int) []int {
	// swap them if min greater than max
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
