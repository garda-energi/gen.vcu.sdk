package sdk

import (
	"errors"
	"log"
)

type Sdk struct {
	logger  *log.Logger
	sleeper Sleeper
	client  *client
}

// New create new instance of Sdk for VCU (Vehicle Control Unit).
func New(cc ClientConfig, logging bool) Sdk {
	logger := newLogger(logging, "SDK")
	return Sdk{
		logger:  logger,
		sleeper: &realSleeper{},
		client:  newClient(&cc, logger),
	}
}

// Connect open connection to mqtt client
func (s *Sdk) Connect() error {
	token := s.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Disconnect close connection to mqtt client
func (s *Sdk) Disconnect() {
	s.client.Disconnect(100)
}

// NewCommander create new instance of commander for specific VIN.
func (s *Sdk) NewCommander(vin int) (*commander, error) {
	return newCommander(vin, s.client, s.sleeper, s.logger)
}

// AddListener subscribe to Status & Report topic (if callback is specified) for spesific vin in range.
//
// Examples :
//
// listen by list :
// s.AddListener(listerner, []int{1, 2 ,3}...)
//
// listen one spesific vin 2341 :
// s.AddListener(listerner, 2341)
//
// listen by range :
// s.AddListener(listerner, sdk.VinRange(min, max)...)
func (s *Sdk) AddListener(ls Listener, vins ...int) error {
	if len(vins) == 0 {
		return errors.New("at least 1 vin supplied")
	}
	if ls.StatusFunc == nil && ls.ReportFunc == nil {
		return errors.New("at least 1 listener supplied")
	}
	if !s.client.IsConnected() {
		return errClientDisconnected
	}

	ls.logger = s.logger
	if ls.StatusFunc != nil {
		topics := setTopicVins(TOPIC_STATUS, vins)
		if err := s.client.subMulti(topics, QOS_SUB_STATUS, ls.status()); err != nil {
			return err
		}
	}

	if ls.ReportFunc != nil {
		topics := setTopicVins(TOPIC_REPORT, vins)
		if err := s.client.subMulti(topics, QOS_SUB_REPORT, ls.report()); err != nil {
			return err
		}
	}
	return nil
}

// RemoveListener unsubscribe status and report topic for spesific vin in range.
func (s *Sdk) RemoveListener(vins ...int) error {
	topics := append(
		setTopicVins(TOPIC_STATUS, vins),
		setTopicVins(TOPIC_REPORT, vins)...,
	)
	return s.client.unsub(topics)
}

// VinRange generate array of integer from min to max.
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
