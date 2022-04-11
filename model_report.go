package sdk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
		if _, ok := data.([]PacketData); ok {
			var idx int
			var err error
			var k_len int = len(k)

			if k_len > 2 && k[0] == '[' && k[k_len-1] == ']' {
				idx, err = strconv.Atoi(k[1 : k_len-1])
				if err != nil {
					return nil
				}
				data = data.([]PacketData)[idx]
				continue
			}
			break
		} else if _, ok := data.(PacketData); !ok {
			break
		}

		data = data.(PacketData)[k]

		if data == nil {
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
		if tag.Tipe == Array_t {
			tag = tag.Sub[0]
			continue
		}
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

// GetBikeError get error from bike's report packet.
// Return BIKE_NOERROR if no error detected
// and return BIKE_ERROR_*** if any error.
// r.GetBikeError().Error() will return error string
func (r *ReportPacket) GetBikeError() BikeError {
	if r.Header.Version == 1 || r.Header.Version == 2 {
		// varsi 1 mcu error
		mcuFaultPost, isFaultPostOK := r.GetValue("Mcu.Faults.Post").(uint32)
		mcuFaultRun, isFaultRunOK := r.GetValue("Mcu.Faults.Run").(uint32)

		if isFaultPostOK && mcuFaultPost > 0 {
			postOverTemp := uint32((1 << MCU_POST_MOD_TEMP_LOW) |
				(1 << MCU_POST_MOD_TEMP_HIGH) |
				(1 << MCU_POST_PCB_TEMP_HIGH) |
				(1 << MCU_POST_GATE_TEMP_HIGH))

			if mcuFaultPost&postOverTemp != 0 {
				return BIKE_ERROR_MOTOR_OVER_TEMPERATURE
			} else {
				return BIKE_ERROR_UNKNOWN
			}

		} else if isFaultRunOK && mcuFaultRun > 0 {
			runOverCurrent := uint32(1 << MCU_RUN_OVER_CURRENT)
			runOverVoltage := uint32(1 << MCU_RUN_OVER_VOLTAGE)
			runUnderVoltage := uint32(1 << MCU_RUN_UNDER_VOLTAGE)
			runOverTemp := uint32((1 << MCU_RUN_INV_OVER_TEMP) |
				(1 << MCU_RUN_MOTOR_OVER_TEMP) |
				(1 << MCU_RUN_MODA_OVER_TEMP) |
				(1 << MCU_RUN_MODB_OVER_TEMP) |
				(1 << MCU_RUN_MODC_OVER_TEMP) |
				(1 << MCU_RUN_PCB_OVER_TEMP) |
				(1 << MCU_RUN_GATE1_OVER_TEMP) |
				(1 << MCU_RUN_GATE2_OVER_TEMP) |
				(1 << MCU_RUN_GATE3_OVER_TEMP))

			if mcuFaultRun&runOverCurrent != 0 {
				return BIKE_ERROR_MOTOR_OVER_CURRENT
			} else if mcuFaultRun&runOverVoltage != 0 {
				return BIKE_ERROR_MOTOR_OVER_VOLTAGE
			} else if mcuFaultRun&runUnderVoltage != 0 {
				return BIKE_ERROR_MOTOR_UNDER_VOLTAGE
			} else if mcuFaultRun&runOverTemp != 0 {
				return BIKE_ERROR_MOTOR_OVER_TEMPERATURE
			} else {
				return BIKE_ERROR_UNKNOWN
			}
		}
	}

	return BIKE_NOERROR
}

// GetBatteryError get error from bike's report packet.
// Return BIKE_BATTERY_NOERROR if no error detected
// and return BIKE_ERROR_BATTERY_*** if any error.
// r.GetBatteryError().Error() will return error string
func (r *ReportPacket) GetBatteryError() BatteryError {
	if r.Header.Version == 1 || r.Header.Version == 2 {
		// Check BMS error
		BmsFault := r.GetValue("Bms.Faults").(uint16)
		if BmsFault > 0 {
			bmsSysFailure := uint16((1 << BMS_SYSTEM_FAILURE))
			bmsOverCurrent := uint16((1 << BMS_DISCHARGE_OVER_CURRENT) | (1 << BMS_CHARGE_OVER_CURRENT))
			bmsOverVoltage := uint16((1 << BMS_OVER_VOLTAGE))
			bmsUnderVoltage := uint16((1 << BMS_UNDER_VOLTAGE))
			bmsOverTemp := uint16((1 << BMS_DISCHARGE_OVER_TEMPERATURE) | (1 << BMS_CHARGE_OVER_TEMPERATURE))
			bmsOverDischarge := uint16((1 << BMS_OVER_DISCHARGE_CAPACITY))

			if BmsFault&bmsSysFailure != 0 {
				return BIKE_ERROR_BATTERY_SYSTEM_FAILURE
			} else if BmsFault&bmsOverCurrent != 0 {
				return BIKE_ERROR_BATTERY_OVER_CURRENT
			} else if BmsFault&bmsOverVoltage != 0 {
				return BIKE_ERROR_BATTERY_OVER_VOLTAGE
			} else if BmsFault&bmsUnderVoltage != 0 {
				return BIKE_ERROR_BATTERY_UNDER_VOLTAGE
			} else if BmsFault&bmsOverTemp != 0 {
				return BIKE_ERROR_BATTERY_OVER_TEMPERATURE
			} else if BmsFault&bmsOverDischarge != 0 {
				return BIKE_ERROR_BATTERY_OVER_DISCHARGE
			} else {
				return BIKE_ERROR_BATTERY_UNKNOWN
			}
		}
	}

	return BIKE_BATTERY_NOERROR
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

// RealtimeData check if current report log is realtime
func (r *ReportPacket) RealtimeData() bool {
	queued, isOK := r.GetValue("Report.Queued").(uint8)
	if !isOK {
		return false
	}
	return queued <= REPORT_REALTIME_QUEUED
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
		if stackVal, isGetStack := stack.(uint8); isGetStack && stackVal < STACK_OVERFLOW_BYTE_MIN {
			return true
		}
	}
	return false
}
