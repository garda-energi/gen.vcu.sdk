package sdk

type Sdk struct {
	logging bool
	broker  Broker
}

// New create new instance of Sdk for VCU (Vehicle Control Unit).
func New(brokerConfig BrokerConfig, logging bool) Sdk {
	return Sdk{
		logging: logging,
		broker:  newBroker(&brokerConfig, logging),
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
	return newCommander(vin, s.broker, &realSleeper{}, s.logging)
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
	ls.logger = newLogger(s.logging, "LISTENER")

	if ls.StatusFunc != nil {
		topics := setTopicToVins(TOPIC_STATUS, vins)
		if err := s.broker.subMulti(topics, QOS_SUB_STATUS, ls.status()); err != nil {
			return err
		}
	}

	if ls.ReportFunc != nil {
		topics := setTopicToVins(TOPIC_REPORT, vins)
		if err := s.broker.subMulti(topics, QOS_SUB_REPORT, ls.report()); err != nil {
			return err
		}
	}
	return nil
}

// RemoveListener unsubscribe status and report topic for spesific vin in range.
func (s *Sdk) RemoveListener(vins ...int) error {
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
