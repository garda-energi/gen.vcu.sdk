package command

import "time"

const FINGERPRINT_MAX = 5
const SPEED_MAX = 110

const (
	DEFAULT_CMD_TIMEOUT = 10 * time.Second
	DEFAULT_ACK_TIMEOUT = 3 * time.Second
)
