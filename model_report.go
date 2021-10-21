package sdk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type ReportPacket struct {
	Header
	Payload message
	Data    PacketData
}

// ValidPrefix check if r's prefix is valid
func (r *ReportPacket) ValidPrefix() bool {
	return r.Header.Prefix == PREFIX_REPORT
}

// ValidSize check if r's size is valid
func (r *ReportPacket) ValidSize(tag tagger) bool {
	// "tagger size" + "vin" + "send datetime"
	return int(r.Header.Size) == tag.getSize()+4+7
}

// GetValue get report packet data by key
func (r *ReportPacket) GetValue(key string) interface{} {
	var data interface{}
	keys := strings.Split(key, ".")
	data = r.Data

	for _, k := range keys {
		data = data.(PacketData)[k]
		if data == nil {
			break
		}
		if _, ok := data.(PacketData); !ok {
			break
		}
	}
	return data
}

// GetValue get report packet data type by key. return VarDataType.
func (r *ReportPacket) GetType(key string) VarDataType {
	var result VarDataType = ""
	rpStructure, isGot := ReportPacketStructures[int(r.Header.Version)]

	if !isGot {
		return result
	}

	tag := rpStructure
	isFound := false

	keys := strings.Split(key, ".")
	for _, k := range keys {
		isFound = false
		for _, subtag := range tag.Sub {
			if subtag.Name == k {
				isFound = true
				tag = subtag
				break
			}
		}
		if !isFound {
			break
		}
	}
	if isFound {
		result = tag.Tipe
	}
	return VarDataType(result)
}

// stringOfData get PacketData as pretty string
func (r *ReportPacket) stringOfData(data interface{}, tag tagger, deep int) string {
	if data == nil {
		return ""
	}

	str := ""
	for i := 0; i < deep; i++ {
		str += " "
	}

	switch tag.Tipe {
	case Struct_t:
		d := data.(PacketData)
		if tag.Name != "" {
			str += fmt.Sprintln(tag.Name, ":")
		} else {
			str = ""
		}
		for _, subtag := range tag.Sub {
			str += r.stringOfData(d[subtag.Name], subtag, deep+1)
		}
	case Array_t:
		s := reflect.ValueOf(data)
		str += fmt.Sprintln(tag.Name, ":")
		for i := 0; i < s.Len(); i++ {
			for i := 0; i < deep+1; i++ {
				str += " "
			}
			str += fmt.Sprintln("#", i, ":")
			for _, subtag := range tag.Sub {
				str += r.stringOfData(s.Index(i).Interface(), subtag, deep+1)
			}
		}
	default:
		str += fmt.Sprintln(tag.Name, ":", data)
	}
	return str
}

// String get report packet data as pretty string
func (r *ReportPacket) String() string {
	rpStructure := ReportPacketStructures[int(r.Header.Version)]
	str := r.stringOfData(r.Data, rpStructure, -1)
	return str
}

// Json get report packet data as json in bytes
func (r *ReportPacket) Json() []byte {
	jsonString, err := json.Marshal(r.Data)
	if err != nil {
		fmt.Println(err)
	}
	return jsonString
}

// VcuGetEvents parse v's events in an array
func (r *ReportPacket) VcuGetEvents() VcuEvents {
	ve := make(VcuEvents, 0, VCU_EVENTS_MAX)
	for i := 0; i < int(VCU_EVENTS_MAX); i++ {
		if r.VcuIsEvents(VcuEvent(i)) {
			ve = append(ve, VcuEvent(i))
		}
	}
	return ve
}

// VcuIsEvents check if v's event has evs
func (r *ReportPacket) VcuIsEvents(ev ...VcuEvent) bool {
	set := 0
	events, isOK := r.GetValue("Vcu.Events").(uint16)
	if !isOK {
		return false
	}

	for _, e := range ev {
		if bitSet(uint32(events), uint8(e)) {
			set++
		}
	}
	return set == len(ev)
}

// VcuRealtimeData check if current report log is realtime
func (r *ReportPacket) VcuRealtimeData() bool {
	logBuffered, isOK := r.GetValue("Vcu.LogBuffered").(uint8)
	if !isOK {
		return false
	}
	return logBuffered <= REPORT_REALTIME_LOG
}

// VcuBatteryLow check if v's backup battery voltage is low
func (r *ReportPacket) VcuBatteryLow() bool {
	batVolt, isOK := r.GetValue("Vcu.BatVoltage").(float32)
	if !isOK {
		return false
	}
	return batVolt < BATTERY_BACKUP_LOW_MV
}

// EepromCapacityLow check if e's storage capacity is low
func (r *ReportPacket) EepromCapacityLow() bool {
	used, isOK := r.GetValue("Eeprom.Used").(uint8)
	if !isOK {
		return false
	}
	return used > EEPROM_LOW_CAPACITY_PERCENT
}

// GpsValidHorizontal check if g's horizontal section (heading, longitude, latitude) is valid
func (r *ReportPacket) GpsValidHorizontal() bool {
	hdop, isHdopOK := r.GetValue("Gps.HDOP").(float32)
	if !isHdopOK {
		return false
	}
	if hdop > GPS_DOP_MIN {
		return false
	}
	longitude, isLongOK := r.GetValue("Gps.Longitude").(float32)
	latitude, isLatOK := r.GetValue("Gps.Latitude").(float32)
	if !(isLongOK && isLatOK) {
		return false
	}
	return (longitude >= GPS_LNG_MIN && longitude <= GPS_LNG_MAX) &&
		(latitude >= GPS_LAT_MIN && latitude <= GPS_LAT_MAX)
}

// GpsValidVertical check if g's vertical section (altitude) is valid
func (r *ReportPacket) GpsValidVertical() bool {
	vdop, isOk := r.GetValue("Gps.VDOP").(float32)
	if !isOk {
		return false
	}
	return vdop <= GPS_DOP_MIN
}

// NetLowSignal check if n's signal is low
func (r *ReportPacket) NetLowSignal() bool {
	signal, isOk := r.GetValue("Net.Signal").(uint8)
	if !isOk {
		return false
	}
	return signal <= NET_LOW_SIGNAL_PERCENT
}

// BmsGetFaults parse b's fault field
func (r *ReportPacket) BmsGetFaults() BmsFaults {
	bf := make(BmsFaults, 0, BMS_FAULTS_MAX)
	for i := 0; i < int(BMS_FAULTS_MAX); i++ {
		if r.BmsIsFaults(BmsFault(i)) {
			bf = append(bf, BmsFault(i))
		}
	}
	return bf
}

// BmsIsFaults check if b's fault has bfs
func (r *ReportPacket) BmsIsFaults(bf ...BmsFault) bool {
	faults, isOk := r.GetValue("Bms.Faults").(uint16)
	if !isOk {
		return false
	}

	set := 0
	for _, f := range bf {
		if bitSet(uint32(faults), uint8(f)) {
			set++
		}
	}
	return set == len(bf)
}

// BmsLowCapacity check if b's SoC (State of Charge) is low
func (r *ReportPacket) BmsLowCapacity() bool {
	soc, isOk := r.GetValue("Bms.SOC").(uint8)
	if !isOk {
		return false
	}

	return soc < BMS_LOW_CAPACITY_PERCENT
}

// McuGetFaults parse mcu's fault field
func (r *ReportPacket) McuGetFaults() McuFaults {
	mf := McuFaults{
		Post: make([]McuFaultPost, 0, MCU_POST_FAULTS_MAX),
		Run:  make([]McuFaultRun, 0, MCU_RUN_FAULTS_MAX),
	}
	for i := 0; i < int(MCU_POST_FAULTS_MAX); i++ {
		if r.McuIsPostFaults(McuFaultPost(i)) {
			mf.Post = append(mf.Post, McuFaultPost(i))
		}
	}
	for i := 0; i < int(MCU_RUN_FAULTS_MAX); i++ {
		if r.McuIsRunFaults(McuFaultRun(i)) {
			mf.Run = append(mf.Run, McuFaultRun(i))
		}
	}
	return mf
}

// McuIsPostFaults check if mcu's post fault has mfs
func (r *ReportPacket) McuIsPostFaults(mf ...McuFaultPost) bool {
	mfp, isOk := r.GetValue("Mcu.Faults.Post").(uint32)
	if !isOk {
		return false
	}

	set := 0
	for _, f := range mf {
		if bitSet(mfp, uint8(f)) {
			set++
		}
	}
	return set == len(mf)
}

// McuIsRunFaults check if mcu's run fault has mfs
func (r *ReportPacket) McuIsRunFaults(mf ...McuFaultRun) bool {
	mfr, isOk := r.GetValue("Mcu.Faults.Run").(uint32)
	if !isOk {
		return false
	}

	set := 0
	for _, f := range mf {
		if bitSet(mfr, uint8(f)) {
			set++
		}
	}
	return set == len(mf)
}

// TaskStackOverflow check if t's stack is near overflow
func (r *ReportPacket) TaskStackOverflow() bool {
	stacks, isOk := r.GetValue("Task.Stack").(PacketData)
	if !isOk {
		return false
	}
	for _, stack := range stacks {
		if stackVal, isGetStack := stack.(uint16); isGetStack && stackVal < STACK_OVERFLOW_BYTE_MIN {
			return true
		}
	}
	return false
}
