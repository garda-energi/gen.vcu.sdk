package sdk

type BikeError byte

const (
	BIKE_NOERROR       BikeError = 0
	BIKE_ERROR_UNKNOWN BikeError = iota
	BIKE_ERROR_MOTOR_OVER_CURRENT
	BIKE_ERROR_MOTOR_OVER_VOLTAGE
	BIKE_ERROR_MOTOR_UNDER_VOLTAGE
	BIKE_ERROR_MOTOR_OVER_TEMPERATURE
	BIKE_ERROR_BATTERY_SYSTEM_FAILURE
	BIKE_ERROR_BATTERY_OVER_CURRENT
	BIKE_ERROR_BATTERY_OVER_VOLTAGE
	BIKE_ERROR_BATTERY_UNDER_VOLTAGE
	BIKE_ERROR_BATTERY_OVER_TEMPERATURE
	BIKE_ERROR_BATTERY_OVER_DISCHARGE
)

var bikeErrorString []string = []string{
	"",
	"Unkown Error",
	"Motor Over Current",
	"Motor Over Voltage",
	"Motor Under Voltage",
	"Motor Over Temperature",
	"Battery System Failure",
	"Battery Over Current",
	"Battery Over Voltage",
	"Battery Under Voltage",
	"Battery Over Temperature",
	"Battery Over Discharge",
}

// Error will return bike error in string
func (errCode BikeError) Error() string {
	return bikeErrorString[errCode]
}
