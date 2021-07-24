package sdk

type Sdk struct {
	broker  Broker
	logging bool
}

// New create new instance of Sdk for VCU (Vehicle Control Unit).
func New(host string, port int, user, pass string, logging bool) Sdk {
	return Sdk{
		broker: newBroker(brokerConfig{
			Host: host,
			Port: port,
			User: user,
			Pass: pass,
		}),
		logging: logging,
	}
}

// Connect open connection to mqtt broker.
func (s *Sdk) Connect() error {
	return s.broker.connect()
}

// Disconnect close connection to mqtt broker.
func (s *Sdk) Disconnect() {
	s.broker.disconnect()
}

// NewCommander create new instance of commander for specific VIN.
func (s *Sdk) NewCommander(vin int) (*commander, error) {
	return newCommander(vin, s.broker, s.logging)
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
		topics := setTopicToVins(TOPIC_STATUS, vins)
		listener := statusListener(l.StatusFunc, s.logging)
		if err := s.broker.subMulti(topics, QOS_SUB_STATUS, listener); err != nil {
			return err
		}
	}

	if l.ReportFunc != nil {
		topics := setTopicToVins(TOPIC_REPORT, vins)
		listener := reportListener(l.ReportFunc, s.logging)
		if err := s.broker.subMulti(topics, QOS_SUB_REPORT, listener); err != nil {
			return err
		}
	}
	return nil
}

// RemoveListener unsubscribe status and report topic for spesific vin in range.
func (s *Sdk) RemoveListener(vins []int) error {
	topics := append(
		setTopicToVins(TOPIC_STATUS, vins),
		setTopicToVins(TOPIC_REPORT, vins)...,
	)
	return s.broker.unsubMulti(topics)
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
